<template>
  <component :is="type" v-bind="linkProps">
    <slot />
  </component>
</template>

<script>
export default {
  name: 'AppLink',
  props: {
    to: {
      type: String,
      required: true
    }
  },
  setup(props) {
    const isExternal = /^(https?:|mailto:|tel:)/.test(props.to)

    const type = computed(() => {
      return isExternal ? 'a' : 'router-link'
    })

    const linkProps = computed(() => {
      if (isExternal) {
        return {
          href: props.to,
          target: '_blank',
          rel: 'noopener'
        }
      }
      return {
        to: props.to
      }
    })

    return {
      type,
      linkProps
    }
  }
}
</script> 