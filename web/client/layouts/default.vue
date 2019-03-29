<template>
  <v-app>
    <Nav/>
    <v-content>
      <v-container fluid fill-height>
        <nuxt/>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
import Nav from "@/components/nav/Navbar"
export default {
  components : {
    Nav,
  }
}
</script>



<style>

</style>
