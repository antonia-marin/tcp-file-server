var vm = new Vue({
    el: '#app',
    data: {
        connection: null,
        channels: [],
        clients: [],
        clientsOnChannels: [],
        files: [],
    },
    methods: {
      requestData: () => {
        setInterval( () => {
          this.connection.send('"/analytics"');;
        }, 5000);
      },
    },
    created: () => {
      console.log("Starting connection to WebSocket Server")
      this.connection = new WebSocket("ws://localhost:8080/ws")

      //var a = this
      this.connection.onmessage = (event) => {
        d = JSON.parse(event.data)
        vm.clients = d["clients"]
        vm.channels = d["channels"]
        vm.clientsOnChannels = d["clientsOnChannels"]
        vm.files = d["files"]
      }

      this.connection.onopen = (event)  => {
        console.log("Successfully connected to the echo websocket server...")
        console.log(event)
      }

      this.connection.onclose = (event)  => {
        console.log("Closes connection: " )
        console.log(event)
      }
    },
    mounted() {
      this.requestData()
    }
});
