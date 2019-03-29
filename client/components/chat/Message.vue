<template >
  <div class="tw-flex tw-mb-6 ">
    <Avatar class="tw-mr-4"/>

    <div class="tw-mt-1 tw-p-4 tw-bg-blue tw-text-white tw-rounded-lg tw-text-sm tw-w-3/5 xl:tw-w-1/3">
      <div class="tw-mb-4">
        {{ data.body }}
        <!-- Lorem ipsum, dolor sit amet consectetur adipisicing elit. Totam nobis, enim fugiat mollitia similique voluptas maxime laboriosam officiis laborum, hic eum error obcaecati, ea magnam dolores voluptatibus. Ab, natus numquam? Lorem ipsum dolor sit amet consectetur adipisicing elit. Distinctio dicta iure at? Mollitia cum dignissimos corrupti nobis dolores est nesciunt non tempore ea nostrum, labore quod quas aspernatur assumenda impedit! -->
      </div>
      <div class="tw-text-right tw-font-sans">
        <span class="tw-mr-2">Sent: </span>
        <span>{{ data.sent }}</span> 
      </div>
    </div>

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
}
</script>

<style>

</style>
