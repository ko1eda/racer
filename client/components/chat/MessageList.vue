<template>
  <div :style="height" class="tw-overflow-y-auto hide-scroll h">
    <Message v-for="(msg,i) in msgs" :key="i" :data="msg"/>
  </div>
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
    height() {
      // for mobile and tablet (subtract less because the margins are smaller 16px vs 24px)
      let h = "calc(100vh - 167px)"

      // anything ipad pro or above will have a padding of 24px so we need to subtract more
      if (this.$vuetify.breakpoint.width >= 1024) h = "calc(100vh - 190px)"


      return {height: h}
    },

    // maxheight() {
    //   let h = "80vh"
      
    //   // portrait mode for phones 
    //   if (this.$vuetify.breakpoint.height < 650) h = "60vh"

    //   if (this.$vuetify.breakpoint.height < 800) h = "74vh"

    //   // ipad
    //   if (this.$vuetify.breakpoint.width < 1024 && this.$vuetify.breakpoint.height >= 1024 ) h = "83vh"
  

    //   // ipad pro
    //   if (this.$vuetify.breakpoint.width === 1024 && this.$vuetify.breakpoint.height === 1366 ) h = "86vh"

    //   return h
    // }
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

