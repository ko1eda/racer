<template>
  <v-responsive :min-height="minheight" :max-height="maxheight" class="tw-overflow-y-auto hide-scroll">
    <Message v-for="(msg,i) in msgs" :key="i" :data="msg"/>
  </v-responsive>
</template>

<script>
import Message from "@/components/chat/Message"
export default {
  components : {
    Message
  },

  props : {
    msgs : {
      type: Array,
      default : () => []
    }
  },

  computed : {
    minheight() {
      let h = "80vh"

      // portrait mode for phones 
      if (this.$vuetify.breakpoint.height < 650) h = "60vh"

      return h
    },

    maxheight() {
      let h = "80vh"
      
      // portrait mode for phones 
      if (this.$vuetify.breakpoint.height < 650) h = "60vh"

      // ipad
      if (this.$vuetify.breakpoint.width < 1024 && this.$vuetify.breakpoint.height >= 1024 ) h = "83vh"
  

      // ipad pro
      if (this.$vuetify.breakpoint.width === 1024 && this.$vuetify.breakpoint.height === 1366 ) h = "86vh"

      return h
    }
  },
}
</script>

<style lang="scss" scoped>
  .hide-scroll::-webkit-scrollbar {
    width: 0px;
    background: transparent; /* Chrome/Safari/Webkit */
  }

  .hide-scroll {
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none;  /* IE 10+ */
  }
</style>

