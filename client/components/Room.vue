<template>
  <div>:_</div>
  <!-- <div class="tw-w-100 tw-flex tw-h-full">
    <div class="tw-w-4/5 tw-flex tw-flex-col">

      <div class="tw-w-full tw-overflow-y-auto tw-flex-1">
        <Message v-for="(message,i) in messages" :key="i" :data="message"/>
      </div>

      <form class="tw-h-16 tw-w-full tw-border tw-flex" @submit.prevent>
        <button class="tw-bg-red tw-text-white tw-h-full tw-w-24 tw-border tw-border-red-light" type="submit" @click="submit">Submit</button>
        <input v-model="form.body" type="text" class="tw-text-sm tw-flex-1 tw-p-4" placeholder="Do it">
      </form>

    </div>

    <div class="tw-flex-1 tw-border">
      user list
    </div>

  </div> -->
</template>

<script>
import Message from "@/components/Message"
export default {
  components : {
    Message
  },

  data() {
    return {
      messages : [],
      conn : WebSocket,
      form : {
        body : ""
      },
      message : {
        senderID : 0,
        sent: "",
        body: "",
      }
    }
  },

  created() {
    // fetch all messages for chatID and scroll screen to bottom
    this.conn = new WebSocket("ws://localhost:80/v1" + this.$route.path);

    this.conn.onmessage = (event) => {
      this.message = JSON.parse(event.data)
      this.messages.push(this.message)
    }
  },


  methods : {
    submit() {
      // create a json string from the message object
      // send it down the pipe 
      this.conn.send(JSON.stringify({body: this.form.body}))

      this.form.body = ""
    }
  }

};
</script>


<style lang="scss" scoped>
  .offset {

  }
</style>