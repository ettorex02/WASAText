<template>
  <div class="d-flex" style="height: 100vh;">
    <!-- Sidebar -->
    <div class="bg-light p-3" style="width: 350px; min-width: 250px; max-width: 420px; border-right: 0; box-sizing: border-box;">
      <!-- Barra di ricerca utenti con dropdown Bootstrap -->
      <div class="dropdown w-100">
        <input
          v-model="search"
          class="form-control mb-3 dropdown-toggle"
          placeholder="Cerca utenti per username..."
          autocomplete="off"
          data-bs-toggle="dropdown"
          @input="searchUsers"
          @focus="dropdownOpen = true"
          @blur="closeDropdown"
        >
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
            <img :src="user.profilePicture" alt="profile" width="32" height="32" class="rounded-circle me-2">
            <span>{{ user.username }}</span>
          </li>
        </ul>
      </div>

      <!-- Qui la lista delle conversazioni -->
      <ul class="list-group mt-3">
        <!-- Conversazioni 1:1 -->
        <li
          v-for="conv in conversations"
          :key="'chat-' + conv.id"
          class="list-group-item list-group-item-action d-flex align-items-center"
          :class="{ 'selected-conv': openConversation && openConversation.id === conv.id }"
          style="cursor:pointer;"
          @click="openConv(conv)"
        >
          <img :src="conv.profilePicture" alt="profile" width="40" class="rounded-circle me-2">
          <div class="flex-grow-1">
            <div class="fw-bold">{{ conv.username }}</div>
            <div class="text-muted small">
              {{ conv.lastMessage && conv.lastMessage.length > 12 ? conv.lastMessage.slice(0, 12) + '…' : conv.lastMessage }}
            </div>
          </div>
          <div class="text-end small text-muted ms-2">
            {{ conv.lastMessageTime }}
          </div>
        </li>
                
        <!-- Gruppi -->
        <li
          v-for="group in groups"
          :key="'group-' + group.id"
          class="list-group-item list-group-item-action d-flex align-items-center"
          :class="{ 'selected-conv': openConversation && openConversation.id === group.id }"
          style="cursor:pointer;"
          @click="openConv({ ...group, username: group.name, profilePicture: group.photo || 'https://cdn-icons-png.flaticon.com/512/74/74472.png' })"
        >
          <img :src="group.photo || 'https://cdn-icons-png.flaticon.com/512/74/74472.png'" alt="group" width="40" class="rounded-circle me-2">
          <div class="flex-grow-1">
            <div class="fw-bold">👥 {{ group.name }}</div>
            <div class="text-muted small">
              <!-- Mostra anteprima solo se la logica backend la fornisce, altrimenti lascia vuoto -->
              {{ group.lastMessage && group.lastMessage.length > 12 ? group.lastMessage.slice(0, 12) + '…' : group.lastMessage }}
            </div>
          </div>
          <div class="text-end small text-muted ms-2">
            {{ group.lastMessageTime }}
          </div>
        </li>
      </ul>

      <!-- In fondo alla sidebar, subito dopo </ul> -->
      <section class="create-group-section mt-4">
        <hr>
        <button class="btn btn-success w-100" @click="openCreateGroupModal = true">
          + Create Group
        </button>
      </section>

      <!-- Modale per creare gruppo -->
      <GroupModal
        :open-create-group-modal="openCreateGroupModal"
        :edit-group-mode="editGroupMode"
        :edit-group-id="editGroupId"
        @close-create-group-modal="closeCreateGroupModal"
        @refresh-groups="loadAll"
      />
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

      <ErrorMsg v-if="errormsg" :msg="errormsg" />

      <!-- Main content -->
      <div class="flex-grow-1 d-flex flex-column justify-content-center" style="height: 100vh;">
        <div
          v-if="openConversation"
          class="chat-box d-flex flex-column"
        >
          <div class="border-bottom p-3 rounded-top bg-white d-flex align-items-center justify-content-between">
            <h5 class="mb-0">{{ openConversation.username }}</h5>
            <div v-if="isGroup" class="ms-auto d-flex align-items-center">
              <AddMemberButton
                :group-id="openConversation.id"
                :current-members="openConversation.members || []"
                @members-added="handleMembersAdded"
                class="me-1"
              />
              <LeaveGroupButton
                :group-id="openConversation.id"
                @left-group="handleGroupLeft"
              />
            </div>
          </div>
          <!-- Sezione messaggi scrollabile -->
          <div class="flex-grow-1 overflow-auto p-3 messages-area">
            <div
              v-for="msg in messages" :key="msg.id"
              class="d-flex mb-3 align-items-center"
              :class="isMyMessage(msg) ? '' : 'flex-row-reverse'"
            >
              <img :src="msg.sender.profilePicture" alt="profile" width="44" height="44" class="rounded-circle mx-2">
              <div>
                <!-- Mostra il nome SOLO se il messaggio NON è mio -->
                <div
                  v-if="!isMyMessage(msg)"
                  class="fw-bold mb-1"
                  style="font-size: 1rem;"
                >
                  {{ msg.sender.displayName || msg.sender.username }}
                </div>
                <div
                  :class="isMyMessage(msg) ? 'msg-sent' : 'msg-received'"
                  style="font-size: 1.3rem; max-width: 520px; word-break: break-word;"
                >
                  <span v-if="msg.is_forwarded || msg.isForwarded" class="badge bg-warning text-dark mb-1" style="font-size: 0.9rem;">
                    Inoltrato
                  </span>
                  {{ msg.content }}
                  <div class="small text-muted mt-1 d-flex align-items-center">
                    <span>{{ msg.timestamp }}</span>
                    <span v-if="isMyMessage(msg)" :class="getStatusClass(msg.status)" style="margin-left: 8px;">
                      {{ getStatusIcon(msg.status) }}
                    </span>
                    <button v-if="isMyMessage(msg)" class="btn btn-sm btn-link text-danger ms-2" title="Elimina" @click="deleteMessage(msg)">
                      🗑️
                    </button>
                    <button
                      class="btn btn-sm btn-link text-primary ms-2"
                      title="Inoltra"
                      @click="openForwardModal(msg)"
                    >
                      ⏩
                    </button>
                  </div>
                </div>
                <MessageReactions
                  :message="msg"
                  :conversation-id="openConversation.id"
                  @refresh="getConversation(openConversation.id)"
                />
              </div>
            </div>
          </div>
          <!-- Barra invio messaggi SEMPRE visibile in basso -->
          <form class="d-flex border-top p-3 bg-white rounded-bottom" @submit.prevent="sendMessage">
            <input v-model="newMessage" class="form-control me-2" placeholder="Scrivi un messaggio..." required>
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
            <!-- Conversazioni 1:1 -->
            <li
              v-for="conv in conversations"
              :key="'chat-' + conv.id"
              class="list-group-item list-group-item-action"
              style="cursor:pointer;"
              @click="forwardMessage(conv.id)"
            >
              <img :src="conv.profilePicture" alt="profile" width="32" class="rounded-circle me-2">
              {{ conv.username }}
            </li>
                
            <!-- Gruppi -->
            <li
              v-for="group in groups"
              :key="'group-' + group.id"
              class="list-group-item list-group-item-action"
              style="cursor:pointer;"
              @click="forwardMessage(group.id)"
            >
              <img :src="group.photo || 'https://cdn-icons-png.flaticon.com/512/74/74472.png'" alt="group" width="32" class="rounded-circle me-2">
              👥 {{ group.name }}
            </li>
          </ul>
          <button class="btn btn-secondary mt-3" @click="closeForwardModal">Annulla</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import MessageReactions from '@/components/MessageReactions.vue';
import GroupModal from '@/components/GroupModal.vue';
import LeaveGroupButton from '@/components/LeaveGroupButton.vue';
import AddMemberButton from '@/components/AddMemberButton.vue'; // <-- Aggiungi questo import

export default {
  components: { MessageReactions, GroupModal, LeaveGroupButton, AddMemberButton }, // <-- Registra il componente
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
      openConversation: null,
      messages: [],
      messagesPolling: null,
      newMessage: "",
      forwardModalOpen: false,
      forwardMsg: null,
      groups: [],
      openCreateGroupModal: false,
      editGroupMode: false,
      editGroupId: null
    }
  },
  computed: {
    isGroup() {
      return this.openConversation && (
        this.openConversation.name ||
        (this.openConversation.username && this.openConversation.username.startsWith('👥'))
      );
    }
  },
  mounted() {
    this.loadAll();
    this.polling = setInterval(this.loadAll, 6000);
    if (this.$route.query.msg) {
      this.successMsg = this.$route.query.msg;
      this.$router.replace({ path: this.$route.path, query: {} });
    }
  },
  beforeUnmount() {
    clearInterval(this.polling);
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
      try {
        const res = await this.$axios.get(`/search/users?q=${encodeURIComponent(this.search)}`, {
          headers: { Authorization: userId }
        });
        const results = res.data;
        this.searchResults = results.filter(u => u.username !== myUsername);
        this.dropdownOpen = !!this.searchResults.length;
      } catch {
        this.searchResults = [];
        this.dropdownOpen = false;
      }
    },
    closeDropdown() {
      setTimeout(() => { this.dropdownOpen = false; }, 150);
    },
    async startConversation(user) {
      const userId = localStorage.getItem("userId");
      try {
        const res = await this.$axios.post("/conversations", { userId: user.id }, {
          headers: { "Content-Type": "application/json", Authorization: userId }
        });
        const data = res.data;
        if (data.conversationId) {
          this.search = "";
          this.searchResults = [];
          this.dropdownOpen = false;
          await this.loadAll();
          const conv = this.conversations.find(c => c.id === data.conversationId);
          if (conv) {
            this.openConv(conv);
          }
        }
      } catch {}
    },
    async getMyConversations() {
      const userId = localStorage.getItem("userId");
      try {
        const res = await this.$axios.get("/conversations", {
          headers: { Authorization: userId }
        });
        this.conversations = res.data;
        if (!this.openConversation && this.conversations.length > 0) {
          this.openConv(this.conversations[0]);
        }
      } catch {
        this.conversations = [];
      }
    },
    async openConv(conv) {
      this.openConversation = conv;
      this.messages = [];
      await this.getConversation(conv.id);
      await this.markMessagesRead();
      this.startMessagesPolling();
    },
    startMessagesPolling() {
      if (this.messagesPolling) clearInterval(this.messagesPolling);
      if (!this.openConversation) return;
      this.messagesPolling = setInterval(async () => {
        if (this.openConversation) {
          await this.getConversation(this.openConversation.id);
          await this.markMessagesRead();
        }
      }, 1000);
    },
    async getConversation(conversationId) {
      const userId = localStorage.getItem("userId");
      try {
        const res = await this.$axios.get(`/conversations/${conversationId}/messages`, {
          headers: { Authorization: userId }
        });
        const msgs = res.data;
        for (const msg of msgs) {
          try {
            const reactionRes = await this.$axios.get(
              `/conversations/${conversationId}/messages/${msg.id}/reactions`,
              { headers: { Authorization: userId } }
            );
            msg.reactions = reactionRes.data;
          } catch {
            msg.reactions = [];
          }
        }
        this.messages = msgs;
      } catch {
        this.messages = [];
      }
    },
    async sendMessage() {
      if (!this.newMessage.trim() || !this.openConversation) return;
      const userId = localStorage.getItem("userId");
      try {
        await this.$axios.post(`/conversations/${this.openConversation.id}/messages`, {
          content: this.newMessage,
          mediaType: "text",
          isForwarded: false
        }, {
          headers: { "Content-Type": "application/json", Authorization: userId }
        });
        this.newMessage = "";
        await this.getConversation(this.openConversation.id);
      } catch {}
    },
    async markMessagesRead() {
      if (!this.openConversation) return;
      const userId = localStorage.getItem("userId");
      try {
        await this.$axios.patch(`/conversations/${this.openConversation.id}/messages/read`, {}, {
          headers: { Authorization: userId }
        });
      } catch {}
    },
    async deleteMessage(msg) {
      if (!confirm("Sei sicuro di voler eliminare questo messaggio?")) return;
      const userId = localStorage.getItem("userId");
      try {
        await this.$axios.delete(`/conversations/${msg.conversation_id}/messages/${msg.id}`, {
          headers: { Authorization: userId }
        });
        this.messages = this.messages.filter(m => m.id !== msg.id);
      } catch {
        alert("Errore durante l'eliminazione del messaggio.");
      }
    },
    isMyMessage(msg) {
      const myId = localStorage.getItem("userId");
      return String(msg.sender.id) === String(myId);
    },
    getStatusIcon(status) {
      if (status === "read") return "✔✔";
      return "✔";
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
    async forwardMessage(targetConversationId) {
      if (!this.forwardMsg) return;
      const userId = localStorage.getItem("userId");
      try {
        await this.$axios.post(
          `/conversations/${this.forwardMsg.conversation_id}/messages/${this.forwardMsg.id}/forward`,
          { targetConversationId },
          { headers: { "Content-Type": "application/json", Authorization: userId } }
        );
        this.messages = this.messages.filter(m => m.id !== this.forwardMsg.id);
        this.successMsg = "Messaggio inoltrato!";
        setTimeout(() => { this.successMsg = null; }, 1000);
        this.closeForwardModal();
      } catch {
        alert("Errore durante l'inoltro del messaggio.");
      }
    },
    async loadAll() {
      await this.getMyConversations();
      await this.listGroups();
      if (this.openConversation) {
        await this.getConversation(this.openConversation.id);
        await this.markMessagesRead();
      }
    },
    async listGroups() {
      const userId = localStorage.getItem("userId");
      try {
        const res = await this.$axios.get("/groups", {
          headers: { Authorization: userId }
        });
        this.groups = res.data;
      } catch {
        this.groups = [];
      }
    },
    closeCreateGroupModal() {
      this.openCreateGroupModal = false;
      this.editGroupMode = false;
      this.editGroupId = null;
    },
    handleGroupLeft(groupId) {
      this.conversations = this.conversations.filter(conv => conv.id !== groupId);
      if (this.openConversation && this.openConversation.id === groupId) {
        this.openConversation = null;
        this.messages = [];
      }
    },
    handleMembersAdded() { this.loadAll(); } // <-- Facoltativo per refresh
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
.create-group-section {
  padding: 1rem 0 0 0;
}
</style>