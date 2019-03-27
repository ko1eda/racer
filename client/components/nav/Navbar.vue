<template>
  <div>
    <Drawer v-model="drawer"/>
    <v-toolbar clipped-right app flat fixed color="white" class="tw-px-4">
      <v-spacer/>
      <v-toolbar-title class="red--text text--darken-3">RC</v-toolbar-title>
      <v-spacer/>
      <v-toolbar-side-icon class="red--text text--darken-3" @click="drawer = !drawer"/>
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
