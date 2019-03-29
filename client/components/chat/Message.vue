<template >
  <div class="tw-flex tw-items-center tw-mb-8">

    <Avatar class="tw-mr-4"/>

    <div 
      class="tw-p-4 tw-bg-blue tw-rounded-lg tw-shadow tw-text-white tw-text-sm tw-w-3/5 xl:tw-w-1/3 tw-cursor-pointer" 
      @click="details = !details">

      <div>
        {{ data.body }}
        <!-- Lorem, ipsum dolor sit amet consectetur adipisicing elit. Tenetur laborum nihil, accusantium dicta corporis illum! Eaque, amet! Temporibus, hic asperiores? Eaque amet laborum eligendi hic error reiciendis explicabo modi facilis? Lorem, ipsum dolor sit amet consectetur adipisicing elit. Inventore, officiis nihil veritatis similique praesentium quas natus libero explicabo illo id maxime vero culpa soluta magnam sit mollitia ipsa molestiae alias. -->
      </div>

      <div v-show="details" class="tw-text-right tw-font-sans tw-mt-2">
        <span class="tw-mr-1">Sent:</span>
        <span>{{ data.sent ? data.sent : "2:04 am" }}</span> 
      </div>
    </div><!-- end message bubble -->

  </div>
</template>

<script>
import Avatar from "@/components/user/Avatar"
export default {
  components: {
    Avatar
  },

  props : {
    data : {
      type: Object,
      required: true
    }
  },

  data() {
    // TODO: Add an indicator for isLeft or isRight
    // to display which side the message should be on 
    // and switch the icon side accordingly (you can use flex order for this I believe or a v-if for simplicity)
    return {
      details : false
    }
  }

}
</script>

<style lang="scss" scoped>
  // .b-rad {
  //   border-radius: 0px 12px 12px 10px;
  // }
</style>

