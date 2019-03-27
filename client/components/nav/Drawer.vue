<template>
  <v-navigation-drawer
    :value="value"
    :mobile-break-point="0"
    fixed
    right
    @input="$emit('input', $event)"
  />
</template>

<script>
export default {
  props : {
    value : {
      type: Boolean,
      default: false
    }
  },
}
</script>

<style>

</style>
