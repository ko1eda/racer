<template functional>
  <div class="tw-mb-2 tw-py-2 tw-px-4 tw-bg-blue tw-text-white tw-text-sm tw-rounded-full">
    <span class="tw-text-sm tw-mr-1 tw-font-sans">
      {{ props.data.sent }} :
    </span>
    {{ props.data.body }}
  </div>
</template>

<script>
export default {
  props : {
    data : {
      type: Object,
      required: true
    }
  }
}
</script>

<style>

</style>
