<template>
    <div class="d-flex" style="height: 100vh;">
        <!-- Sidebar -->
        <div class="bg-light p-3" style="width: 350px; min-width: 250px; max-width: 420px; border-right: 0; box-sizing: border-box;">
            <!-- Barra di ricerca utenti con dropdown Bootstrap -->
            <div class="dropdown w-100">
                <input
                    v-model="search"
                    @input="searchUsers"
                    class="form-control mb-3 dropdown-toggle"
                    placeholder="Cerca utenti per username..."
                    autocomplete="off"
                    data-bs-toggle="dropdown"
                    @focus="dropdownOpen = true"
                    @blur="closeDropdown"
                />
                <ul
                    class="dropdown-menu w-100"
                    :class="{ show: dropdownOpen && searchResults.length }"
                    style="max-height: 300px; overflow-y: auto;"
                >
                    <li
                        v-for="user in searchResults"
                        :key="user.username"
                        class="dropdown-item d-flex align-items-center"
                        @mousedown.prevent="startConversation(user)"
                    >
                        <img :src="user.profilePicture" alt="profile" width="32" height="32" class="rounded-circle me-2" />
                        <span>{{ user.username }}</span>
                    </li>
                </ul>
            </div>

            <!-- Qui la lista delle conversazioni -->
            <ul class="list-group mt-3">
                <li
                    v-for="conv in conversations"
                    :key="conv.id"
                    class="list-group-item list-group-item-action d-flex align-items-center"
                    :class="{ 'selected-conv': openConversation && openConversation.id === conv.id }"
                    @click="openConv(conv)"
                    style="cursor:pointer;"
                >
                    <img :src="conv.profilePicture" alt="profile" width="40" class="rounded-circle me-2" />
                    <div class="flex-grow-1">
                        <div class="fw-bold">{{ conv.username }}</div>
                        <div class="text-muted small">
                            {{ conv.lastMessage.length > 12 ? conv.lastMessage.slice(0, 12) + '…' : conv.lastMessage }}
                        </div>
                    </div>
                    <div class="text-end small text-muted ms-2">
                        {{ conv.lastMessageTime }}
                    </div>
                </li>
            </ul>
        </div>

        <!-- Main content -->
        <div class="flex-grow-1">
            <div style="position: absolute; top: 80px; right: 30px; z-index: 10;">
                <button class="btn btn-outline-primary me-2" @click="goToProfile">Profilo</button>
                <button class="btn btn-danger" @click="logout">Logout</button>
            </div>

            <div v-if="successMsg" class="alert alert-success text-center" style="max-width: 500px; margin: 0 auto 20px auto;">
                {{ successMsg }}
            </div>

            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

            <!-- Main content -->
            <div class="flex-grow-1 d-flex flex-column justify-content-center" style="height: 100vh;">
                <div
                    v-if="openConversation"
                    class="chat-box d-flex flex-column"
                >
                    <div class="border-bottom p-3 rounded-top bg-white">
                        <h5 class="mb-0">{{ openConversation.username }}</h5>
                    </div>
                    <!-- Sezione messaggi scrollabile -->
                    <div class="flex-grow-1 overflow-auto p-3 messages-area">
                        <div v-for="msg in messages" :key="msg.id"
                            class="d-flex mb-3 align-items-end"
                            :class="isMyMessage(msg) ? '' : 'flex-row-reverse'">
                            <img :src="msg.sender.profilePicture" alt="profile" width="44" height="44" class="rounded-circle mx-2" />
                            <div>
                                <div :class="isMyMessage(msg) ? 'msg-sent' : 'msg-received'"
                                    style="font-size: 1.3rem; max-width: 520px; word-break: break-word;">
                                    <span v-if="msg.is_forwarded || msg.isForwarded" class="badge bg-warning text-dark mb-1" style="font-size: 0.9rem;">
                                        Inoltrato
                                    </span>
                                    {{ msg.content }}
                                    <div class="small text-muted mt-1 d-flex align-items-center">
                                        <span>{{ msg.timestamp }}</span>
                                        <span v-if="isMyMessage(msg)" :class="getStatusClass(msg.status)" style="margin-left: 8px;">
                                            {{ getStatusIcon(msg.status) }}
                                        </span>
                                        <button v-if="isMyMessage(msg)" @click="deleteMessage(msg)" class="btn btn-sm btn-link text-danger ms-2" title="Elimina">
                                            🗑️
                                        </button>
                                        <button
                                          @click="openForwardModal(msg)"
                                          class="btn btn-sm btn-link text-primary ms-2"
                                          title="Inoltra"
                                        >
                                          ⏩
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <!-- Barra invio messaggi SEMPRE visibile in basso -->
                    <form @submit.prevent="sendMessage" class="d-flex border-top p-3 bg-white rounded-bottom">
                        <input v-model="newMessage" class="form-control me-2" placeholder="Scrivi un messaggio..." required />
                        <button class="btn btn-primary" type="submit">Invia</button>
                    </form>
                </div>
                <div v-else class="flex-grow-1 d-flex align-items-center justify-content-center text-muted">
                    Nessuna conversazione trovata. Inizia a cercare utenti per iniziare una chat!
                </div>
            </div>
        </div>

        <!-- Modal per inoltrare messaggio -->
        <div v-if="forwardModalOpen" class="modal-backdrop">
          <div class="modal-dialog">
            <div class="modal-content p-3">
              <h5>Scegli la chat dove inoltrare</h5>
              <ul class="list-group">
                <li
                  v-for="conv in conversations"
                  :key="conv.id"
                  class="list-group-item list-group-item-action"
                  @click="forwardMessageTo(conv.id)"
                  style="cursor:pointer;"
                >
                  <img :src="conv.profilePicture" alt="profile" width="32" class="rounded-circle me-2" />
                  {{ conv.username }}
                </li>
              </ul>
              <button class="btn btn-secondary mt-3" @click="closeForwardModal">Annulla</button>
            </div>
          </div>
        </div>
    </div>
</template>

<script>
export default {
    data() {
        return {
            errormsg: null,
            loading: false,
            some_data: null,
            successMsg: null,
            search: "",
            searchResults: [],
            dropdownOpen: false,
            openConversationId: null,
            conversations: [],
            polling: null,
            openConversation: null, // oggetto conversazione selezionata
            messages: [],
            messagesPolling: null,
            newMessage: "",
            forwardModalOpen: false,
            forwardMsg: null,
        }
    },
    methods: {
        goToProfile() {
            this.$router.push('/profile');
        },
        logout() {
            localStorage.clear();
            this.$router.push('/');
        },
        async refresh() {
            this.loading = true;
            this.errormsg = null;
            try {
                let response = await this.$axios.get("/");
                this.some_data = response.data;
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },
        async searchUsers() {
            if (this.search.length < 1) {
                this.searchResults = [];
                this.dropdownOpen = false;
                return;
            }
            const userId = localStorage.getItem("userId");
            const myUsername = localStorage.getItem("username");
            const res = await fetch(`${__API_URL__}/search/users?q=${encodeURIComponent(this.search)}`, {
                headers: { Authorization: userId }
            });
            if (res.ok) {
                const results = await res.json();
                this.searchResults = results.filter(u => u.username !== myUsername);
                this.dropdownOpen = !!this.searchResults.length;
            } else {
                this.searchResults = [];
                this.dropdownOpen = false;
            }
        },
        closeDropdown() {
            setTimeout(() => { this.dropdownOpen = false; }, 150);
        },
        async startConversation(user) {
            const userId = localStorage.getItem("userId"); // chi ha cercato
            console.log("Utente loggato (chi cerca):", userId);
            console.log("Utente selezionato dalla ricerca:", user);

            const res = await fetch(`${__API_URL__}/conversations`, {
                method: "POST",
                headers: { "Content-Type": "application/json", Authorization: userId },
                body: JSON.stringify({ userId: user.id })
            });
            const data = await res.json();
            console.log("Risposta backend creazione conversazione:", data);

            if (res.ok && data.conversationId) {
                this.search = "";
                this.searchResults = [];
                this.dropdownOpen = false;
                await this.loadConversations(); // aggiorna subito la lista
            } else {
                // Rimuovo l'alert fastidioso, mostro solo in console
                console.log("Errore nella creazione della conversazione:", data.message || data);
            }
        },
        async loadConversations() {
            const userId = localStorage.getItem("userId");
            const res = await fetch(`${__API_URL__}/conversations`, {
                headers: { Authorization: userId }
            });
            if (res.ok) {
                this.conversations = await res.json();
                // Se non c'è una conversazione aperta, apri la prima
                if (!this.openConversation && this.conversations.length > 0) {
                    this.openConv(this.conversations[0]);
                }
            } else {
                this.conversations = [];
            }
        },
        async openConv(conv) {
            this.openConversation = conv;
            await this.loadMessages(conv.id);
            await this.markMessagesRead(); // Segna come letti subito se la chat è aperta
            if (this.messagesPolling) clearInterval(this.messagesPolling);
            this.messagesPolling = setInterval(async () => {
                if (this.openConversation) {
                    await this.loadMessages(this.openConversation.id);
                    await this.markMessagesRead(); // Segna come letti ogni volta che arrivano nuovi messaggi e la chat è aperta
                }
            }, 1000);
        },
        async loadMessages(conversationId) {
            const userId = localStorage.getItem("userId");
            const res = await fetch(`${__API_URL__}/conversations/${conversationId}/messages`, {
                headers: { Authorization: userId }
            });
            if (res.ok) {
                this.messages = await res.json();
            } else {
                this.messages = [];
            }
        },
        async sendMessage() {
            if (!this.newMessage.trim() || !this.openConversation) return;
            const userId = localStorage.getItem("userId");
            const res = await fetch(`${__API_URL__}/conversations/${this.openConversation.id}/messages`, {
                method: "POST",
                headers: { "Content-Type": "application/json", Authorization: userId },
                body: JSON.stringify({ content: this.newMessage, mediaType: "text", isForwarded: false })
            });
            if (res.ok) {
                this.newMessage = "";
                // Non aggiornare qui i messaggi, ci pensa il polling!
            }
        },
        async markMessagesRead() {
            if (!this.openConversation) return;
            const userId = localStorage.getItem("userId");
            await fetch(`${__API_URL__}/conversations/${this.openConversation.id}/messages/read`, {
                method: "PATCH",
                headers: { Authorization: userId }
            });
        },
        async deleteMessage(msg) {
            if (!confirm("Sei sicuro di voler eliminare questo messaggio?")) return;
            const userId = localStorage.getItem("userId");
            const res = await fetch(`${__API_URL__}/conversations/${msg.conversation_id}/messages/${msg.id}`, {
                method: "DELETE",
                headers: { Authorization: userId }
            });
            if (res.ok) {
                // Aggiorna la lista (il polling farà il resto, ma puoi anche rimuovere subito il messaggio)
                this.messages = this.messages.filter(m => m.id !== msg.id);
            } else {
                alert("Errore durante l'eliminazione del messaggio.");
            }
        },
        isMyMessage(msg) {
            const myId = localStorage.getItem("userId");
            // Confronta come stringhe per sicurezza
            return String(msg.sender.id) === String(myId);
        },
        getStatusIcon(status) {
            if (status === "read") return "✔✔✔";      // Tre spunte blu
            return "✔✔";                              // Due spunte grigie per tutto il resto
        },
        getStatusClass(status) {
            if (status === "read") return "text-primary";
            return "text-secondary";
        },
        openForwardModal(msg) {
          this.forwardMsg = msg;
          this.forwardModalOpen = true;
        },
        closeForwardModal() {
          this.forwardMsg = null;
          this.forwardModalOpen = false;
        },
        async forwardMessageTo(targetConversationId) {
          if (!this.forwardMsg) return;
          const userId = localStorage.getItem("userId");
          const res = await fetch(
            `${__API_URL__}/conversations/${this.forwardMsg.conversation_id}/messages/${this.forwardMsg.id}/forward`,
            {
              method: "POST",
              headers: { "Content-Type": "application/json", Authorization: userId },
              body: JSON.stringify({ targetConversationId }),
            }
          );
          if (res.ok) {
            this.successMsg = "Messaggio inoltrato!";
            this.closeForwardModal();
          } else {
            alert("Errore durante l'inoltro del messaggio.");
          }
        }
    },
    mounted() {
        this.loadConversations();
        this.polling = setInterval(this.loadConversations, 1500); // polling ogni 1.5 secondi
        if (this.$route.query.msg) {
            this.successMsg = this.$route.query.msg;
            this.$router.replace({ path: this.$route.path, query: {} });
        }
    },
    beforeUnmount() {
        clearInterval(this.polling);
        if (this.messagesPolling) clearInterval(this.messagesPolling);
    }
}
</script>

<style>
.selected-conv {
    background: rgba(0, 123, 255, 0.15) !important;
}

.chat-box {
    width: 100%;
    height: 90vh;
    min-height: 500px;
    background: #f8f9fa;
    border-radius: 24px;
    box-shadow: 0 2px 16px rgba(0,0,0,0.08);
    overflow: hidden;
    display: flex;
    flex-direction: column;
}
.messages-area {
    flex-grow: 1;
    min-height: 0;
    overflow-y: auto;
    background: #f8f9fa;
}
.msg-sent {
    background: #9accff;    /* sfondo chiaro per i tuoi messaggi */
    color: #222;
    display: inline-block;
    padding: 1rem;
    border-radius: 1.5rem;
}
.msg-received {
    background: #67b3ff;    /* sfondo scuro per i messaggi dell'altro utente */
    color: #222;
    display: inline-block;
    padding: 1rem;
    border-radius: 1.5rem;
}
.modal-backdrop {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.3);
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
}
.modal-dialog {
  background: #fff;
  border-radius: 12px;
  max-width: 400px;
  width: 100%;
}
</style>
