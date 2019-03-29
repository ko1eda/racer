<template>
  <div>
    <Drawer v-model="drawer"/>
    <v-toolbar clipped-right app flat fixed color="white" class="tw-px-4">
      <v-icon color="red" class="tw-cursor-pointer">arrow_back</v-icon>
      <v-spacer/>
      <v-toolbar-title class="red--text text--darken-2">RC</v-toolbar-title>
      <v-spacer/>
      <v-toolbar-side-icon class="red--text text--darken-2" @click="drawer = !drawer"/>
    </v-toolbar>
  </div>
</template>

<script>
import Drawer from "@/components/nav/Drawer"
export default {
   components : {
    Drawer
  },

  data() {
    return {
      drawer : false
    }
  }

}
</script>

<style>

</style>
