<template>
  <div class="tw-flex tw-flex-col tw-w-100 tw-h-screen">
    <div class="tw-flex tw-h-full">

      <div class="tw-w-4/5 tw-p-4 tw-overflow-y-auto tw-h-full tw-border-r tw-border-grey-light">
        <Message v-for="(message,i) in messages" :key="i" :data="message"/>
      </div>

      <div class="tw">
        user list
      </div>

    </div>

    <form class="tw-h-32 tw-w-full tw-border-black tw-border-8 tw-flex" @submit.prevent >
      <button class="tw-bg-red tw-text-white tw-h-full tw-w-24 tw-border tw-border-red-light" type="submit" @click="submit">Submit</button>
      <textarea v-model="form.body" type="text" class="tw-text-sm tw-flex-1 tw-border-none tw-p-4 " placeholder="Do it"/>
    </form>

  </div>
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
    this.conn = new WebSocket("ws://localhost:80/racer" + this.$route.path);

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

<style>

</style>
