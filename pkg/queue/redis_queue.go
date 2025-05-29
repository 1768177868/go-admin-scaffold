package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisQueue Redis队列驱动
type RedisQueue struct {
	client *redis.Client
	config Config
}

// NewRedisQueue 创建Redis队列驱动
func NewRedisQueue(config Config) (*RedisQueue, error) {
	// 解析Redis连接配置
	connection, ok := config.Options["connection"].(string)
	if !ok {
		return nil, fmt.Errorf("redis connection string not found in options")
	}

	opt, err := redis.ParseURL(connection)
	if err != nil {
		return nil, fmt.Errorf("invalid redis connection string: %v", err)
	}

	client := redis.NewClient(opt)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return &RedisQueue{
		client: client,
		config: config,
	}, nil
}

// Push 推送任务到队列
func (q *RedisQueue) Push(ctx context.Context, job JobInterface) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	queue := job.GetQueue()
	if queue == "" {
		if queueOpt, ok := q.config.Options["queue"].(string); ok {
			queue = queueOpt
		} else {
			return fmt.Errorf("queue name not found in options")
		}
	}

	// 如果有延迟，使用ZADD
	if job.GetDelay() > 0 {
		score := float64(time.Now().Add(job.GetDelay()).Unix())
		return q.client.ZAdd(ctx, fmt.Sprintf("queues:%s:delayed", queue), redis.Z{
			Score:  score,
			Member: payload,
		}).Err()
	}

	// 否则直接推送到队列
	return q.client.LPush(ctx, fmt.Sprintf("queues:%s", queue), payload).Err()
}

// PushRaw 推送原始数据到队列
func (q *RedisQueue) PushRaw(ctx context.Context, queue string, payload []byte, options map[string]interface{}) error {
	if queue == "" {
		if queueOpt, ok := q.config.Options["queue"].(string); ok {
			queue = queueOpt
		} else {
			return fmt.Errorf("queue name not found in options")
		}
	}

	delay := time.Duration(0)
	if v, ok := options["delay"].(time.Duration); ok {
		delay = v
	}

	if delay > 0 {
		score := float64(time.Now().Add(delay).Unix())
		return q.client.ZAdd(ctx, fmt.Sprintf("queues:%s:delayed", queue), redis.Z{
			Score:  score,
			Member: payload,
		}).Err()
	}

	return q.client.LPush(ctx, fmt.Sprintf("queues:%s", queue), payload).Err()
}

// Later 延迟推送任务
func (q *RedisQueue) Later(ctx context.Context, job JobInterface, delay time.Duration) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	queue := job.GetQueue()
	if queue == "" {
		if queueOpt, ok := q.config.Options["queue"].(string); ok {
			queue = queueOpt
		} else {
			return fmt.Errorf("queue name not found in options")
		}
	}

	score := float64(time.Now().Add(delay).Unix())
	return q.client.ZAdd(ctx, fmt.Sprintf("queues:%s:delayed", queue), redis.Z{
		Score:  score,
		Member: payload,
	}).Err()
}

// Pop 从队列中取出任务
func (q *RedisQueue) Pop(ctx context.Context, queue string) (JobInterface, error) {
	if queue == "" {
		if queueOpt, ok := q.config.Options["queue"].(string); ok {
			queue = queueOpt
		} else {
			return nil, fmt.Errorf("queue name not found in options")
		}
	}

	// 检查延迟队列
	now := float64(time.Now().Unix())
	delayedKey := fmt.Sprintf("queues:%s:delayed", queue)

	// 获取到期的延迟任务
	results, err := q.client.ZRangeByScore(ctx, delayedKey, &redis.ZRangeBy{
		Min:    "0",
		Max:    fmt.Sprintf("%f", now),
		Offset: 0,
		Count:  1,
	}).Result()

	if err != nil {
		return nil, err
	}

	// 如果有到期的延迟任务，移除并返回
	if len(results) > 0 {
		payload := results[0]
		q.client.ZRem(ctx, delayedKey, payload)

		var job BaseJob
		if err := json.Unmarshal([]byte(payload), &job); err != nil {
			return nil, err
		}
		return &job, nil
	}

	// 从主队列获取任务
	payload, err := q.client.RPop(ctx, fmt.Sprintf("queues:%s", queue)).Bytes()
	if err == redis.Nil {
		return nil, ErrQueueEmpty
	}
	if err != nil {
		return nil, err
	}

	var job BaseJob
	if err := json.Unmarshal(payload, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

// Size 获取队列大小
func (q *RedisQueue) Size(ctx context.Context, queue string) (int64, error) {
	if queue == "" {
		if queueOpt, ok := q.config.Options["queue"].(string); ok {
			queue = queueOpt
		} else {
			return 0, fmt.Errorf("queue name not found in options")
		}
	}

	// 获取主队列大小
	size, err := q.client.LLen(ctx, fmt.Sprintf("queues:%s", queue)).Result()
	if err != nil {
		return 0, err
	}

	// 获取延迟队列大小
	delayedSize, err := q.client.ZCard(ctx, fmt.Sprintf("queues:%s:delayed", queue)).Result()
	if err != nil {
		return 0, err
	}

	return size + delayedSize, nil
}

// Delete 删除任务
func (q *RedisQueue) Delete(ctx context.Context, queue string, job JobInterface) error {
	if queue == "" {
		if queueOpt, ok := q.config.Options["queue"].(string); ok {
			queue = queueOpt
		} else {
			return fmt.Errorf("queue name not found in options")
		}
	}

	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	// 从延迟队列中删除
	delayedKey := fmt.Sprintf("queues:%s:delayed", queue)
	q.client.ZRem(ctx, delayedKey, payload)

	// 从主队列中删除
	// 注意：Redis的List不支持直接删除指定元素，这里需要遍历
	// 实际应用中可能需要更复杂的实现
	key := fmt.Sprintf("queues:%s", queue)
	items, err := q.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return err
	}

	pipe := q.client.Pipeline()
	for _, item := range items {
		if item == string(payload) {
			pipe.LRem(ctx, key, 1, item)
		}
	}
	_, err = pipe.Exec(ctx)
	return err
}

// Release 释放任务回队列
func (q *RedisQueue) Release(ctx context.Context, queue string, job JobInterface, delay time.Duration) error {
	if queue == "" {
		if queueOpt, ok := q.config.Options["queue"].(string); ok {
			queue = queueOpt
		} else {
			return fmt.Errorf("queue name not found in options")
		}
	}

	// 增加重试次数
	job.SetAttempts(job.GetAttempts() + 1)

	// 如果超过最大重试次数，直接删除
	if job.GetAttempts() >= job.GetMaxAttempts() {
		return q.Delete(ctx, queue, job)
	}

	// 计算新的延迟时间
	newDelay := delay
	if newDelay == 0 {
		// 使用退避策略
		backoff := job.GetBackoff()
		if len(backoff) > 0 {
			attempt := job.GetAttempts() - 1
			if attempt < len(backoff) {
				newDelay = backoff[attempt]
			} else {
				newDelay = backoff[len(backoff)-1]
			}
		} else {
			newDelay = job.GetRetryAfter()
		}
	}

	// 重新入队
	return q.Later(ctx, job, newDelay)
}

// Clear 清空队列
func (q *RedisQueue) Clear(ctx context.Context, queue string) error {
	if queue == "" {
		if queueOpt, ok := q.config.Options["queue"].(string); ok {
			queue = queueOpt
		} else {
			return fmt.Errorf("queue name not found in options")
		}
	}

	pipe := q.client.Pipeline()
	pipe.Del(ctx, fmt.Sprintf("queues:%s", queue))
	pipe.Del(ctx, fmt.Sprintf("queues:%s:delayed", queue))
	_, err := pipe.Exec(ctx)
	return err
}

// Close 关闭队列连接
func (q *RedisQueue) Close() error {
	return q.client.Close()
}
