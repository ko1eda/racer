<template>
  <div class="tw-mb-2 tw-py-2 tw-px-4 tw-bg-blue tw-text-white tw-text-sm tw-rounded-full">
    {{ body }}
  </div>
</template>

<script>
export default {
  props : {
    body : {
      type: String,
      default : ""
    }
  }
}
</script>

<style>

</style>
