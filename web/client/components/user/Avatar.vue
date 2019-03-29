<template>
  <v-avatar :size="size" color="grey">
    <!-- obviously replace this with the users actual avatar -->
    <img src="http://i.pravatar.cc/300">
  </v-avatar>
</template>

<script>
export default {
  props : {
    size : {
      type: String,
      default : "48px"
    }
  }
}
</script>

<style>

</style>
