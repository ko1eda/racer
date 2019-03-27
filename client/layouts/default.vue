<template>
  <v-app>
    <Nav/>
    <v-content>
      <v-container fluid>
        <nuxt/>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
import Nav from "@/components/nav/Navbar"
// import Drawer from "@/components/nav/Drawer"
export default {
  components : {
    Nav,
    // Drawer
  }
}
</script>



<style>

</style>
