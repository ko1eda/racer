<template>
  <v-layout>
    <v-flex sm10 class="">
      <v-layout column class="tw-h-full">
        <v-flex>
          <MessageList :msgs="messages"/>
        </v-flex>
        <v-flex shrink class="">
          <v-text-field
            box
            append-outer-icon="send"
            clearable
            label="Chat"
            type="text"
            hide-details
            @click:append-outer="submit"
          />
        </v-flex>
      </v-layout>
    </v-flex>
  </v-layout>
</template>

<script>
import MessageList from "@/components/chat/MessageList"
export default {
  components : {
    MessageList
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
  .drawer-offset {
    width: 305px;
  }
</style>




  <!-- <v-container>
    <v-layout column>
      <v-flex xs8>
        <div class="tw-w-full tw-overflow-y-auto tw-flex-1">
          <Message v-for="(message,i) in messages" :key="i" :data="message"/>
        </div>
        <v-text-field
          box
          clear-icon="md-close-circle"
          clearable
          label="Message"
          type="text"
        />
      </v-flex>
    </v-layout>
  </v-container> -->
  <!--  <Message v-for="(message,i) in messages" :key="i" :data="message"/> -->