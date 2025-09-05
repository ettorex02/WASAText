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
                <!-- Conversazioni 1:1 -->
                <li
                    v-for="conv in conversations"
                    :key="'chat-' + conv.id"
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
                
                <!-- Gruppi -->
                <li
                    v-for="group in groups"
                    :key="'group-' + group.id"
                    class="list-group-item list-group-item-action d-flex align-items-center"
                    :class="{ 'selected-conv': openConversation && openConversation.id === group.id }"
                    @click="openConv({...group, username: group.name, profilePicture: group.photo || 'https://cdn-icons-png.flaticon.com/512/74/74472.png'})"
                    style="cursor:pointer;"
                >
                    <img :src="group.photo || 'https://cdn-icons-png.flaticon.com/512/74/74472.png'" alt="group" width="40" class="rounded-circle me-2" />
                    <div class="flex-grow-1">
                        <div class="fw-bold">👥 {{ group.name }}</div>
                        <div class="text-muted small">
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
              <hr />
              <button class="btn btn-success w-100" @click="openCreateGroupModal = true">
                + Create Group
              </button>
            </section>

            <!-- Modale per creare gruppo -->
            <div v-if="openCreateGroupModal" class="modal-backdrop">
              <div class="modal-dialog">
                <div class="modal-content p-4">
                  <h5>{{ editGroupMode ? 'Modifica gruppo' : 'Crea nuovo gruppo' }}</h5>
                  <form @submit.prevent="editGroupMode ? submitEditGroup() : submitGroup()">
                    <div class="mb-3">
                      <label class="form-label">Nome gruppo</label>
                      <input v-model="newGroupName" class="form-control" required minlength="3" maxlength="16" />
                    </div>
                    <div class="mb-3">
                      <label class="form-label">URL foto gruppo (opzionale)</label>
                      <input v-model="newGroupPhoto" class="form-control" placeholder="https://..." />
                    </div>
                    <div class="mb-3">
                      <label class="form-label">Aggiungi membri</label>
                      <!-- Barra di ricerca utenti identica a quella già presente -->
                      <div class="dropdown w-100">
                        <input
                          v-model="searchUser"
                          @input="searchUsersGroup"
                          class="form-control mb-2 dropdown-toggle"
                          placeholder="Cerca utenti per username..."
                          autocomplete="off"
                          data-bs-toggle="dropdown"
                          @focus="dropdownOpenGroup = true"
                          @blur="closeDropdownGroup"
                        />
                        <ul
                          class="dropdown-menu w-100"
                          :class="{ show: dropdownOpenGroup && searchResultsGroup.length }"
                          style="max-height: 200px; overflow-y: auto;"
                        >
                          <li
                            v-for="user in searchResultsGroup"
                            :key="user.username"
                            class="dropdown-item d-flex align-items-center"
                            @mousedown.prevent="addToGroup(user)"
                          >
                            <img :src="user.profilePicture" alt="profile" width="32" height="32" class="rounded-circle me-2" />
                            <span>{{ user.username }}</span>
                          </li>
                        </ul>
                      </div>
                      <div class="mt-2">
                        <span
                          v-for="user in groupMembers"
                          :key="user.username"
                          class="badge bg-primary me-2"
                          style="font-size:1rem;"
                        >
                          <img :src="user.profilePicture" alt="profile" width="20" height="20" class="rounded-circle me-1" />
                          {{ user.username }}
                          <span class="ms-1" style="cursor:pointer;" @click="removeMemberFromGroup(user.username)">×</span>
                        </span>
                      </div>
                    </div>
                    <div class="d-flex justify-content-end gap-2">
                      <button type="button" class="btn btn-secondary" @click="closeCreateGroupModal">Annulla</button>
                      <button type="submit" class="btn btn-primary">
                        {{ editGroupMode ? 'Salva modifiche' : 'Crea' }}
                      </button>
                    </div>
                    <div v-if="groupError" class="alert alert-danger mt-2">{{ groupError }}</div>
                  </form>
                </div>
              </div>
            </div>
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
                    <div class="border-bottom p-3 rounded-top bg-white d-flex align-items-center justify-content-between">
                      <h5 class="mb-0">{{ openConversation.username }}</h5>
                      <div>
                        <button
                          v-if="isGroup"
                          class="btn btn-outline-primary btn-sm me-2"
                          @click="openEditGroupModal"
                        >
                          Edit Group
                        </button>
                        <button
                          v-if="isGroup"
                          class="btn btn-outline-danger btn-sm"
                          @click="leaveGroup(openConversation.id)"
                        >
                          Lascia Gruppo
                        </button>
                      </div>
                    </div>
                    <!-- Sezione messaggi scrollabile -->
                    <div class="flex-grow-1 overflow-auto p-3 messages-area">
                        <div v-for="msg in messages" :key="msg.id"
                            class="d-flex mb-3 align-items-center"
                            :class="isMyMessage(msg) ? '' : 'flex-row-reverse'">
                            <img :src="msg.sender.profilePicture" alt="profile" width="44" height="44" class="rounded-circle mx-2" />
                            <div>
                                <!-- Mostra il nome SOLO se il messaggio NON è mio -->
                                <div
                                  v-if="!isMyMessage(msg)"
                                  class="fw-bold mb-1"
                                  style="font-size: 1rem;">
                                  {{ msg.sender.displayName || msg.sender.username }}
                                </div>
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
                                <MessageReactions
                                  :message="msg"
                                  :conversationId="openConversation.id"
                                  @refresh="getConversation(openConversation.id)"
                                />
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
                <!-- Conversazioni 1:1 -->
                <li
                  v-for="conv in conversations"
                  :key="'chat-' + conv.id"
                  class="list-group-item list-group-item-action"
                  @click="forwardMessage(conv.id)"
                  style="cursor:pointer;"
                >
                  <img :src="conv.profilePicture" alt="profile" width="32" class="rounded-circle me-2" />
                  {{ conv.username }}
                </li>
                
                <!-- Gruppi -->
                <li
                  v-for="group in groups"
                  :key="'group-' + group.id"
                  class="list-group-item list-group-item-action"
                  @click="forwardMessage(group.id)"
                  style="cursor:pointer;"
                >
                  <img :src="group.photo || 'https://cdn-icons-png.flaticon.com/512/74/74472.png'" alt="group" width="32" class="rounded-circle me-2" />
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

export default {
  components: { MessageReactions },
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
            newGroupName: "",
            newGroupPhoto: "",
            groupMembers: [],
            groupError: "",
            searchUser: "",
            searchResultsGroup: [],
            dropdownOpenGroup: false,
            groups: [],
            openCreateGroupModal: false, 
            editGroupMode: false,
            editGroupId: null,
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
            const userId = localStorage.getItem("userId");
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
                await this.loadAll(); // aggiorna la lista unificata
                
                // Apri automaticamente la conversazione appena creata/trovata
                const conv = this.conversations.find(c => c.id === data.conversationId);
                if (conv) {
                    this.openConv(conv);
                }
            } else {
                console.log("Errore nella creazione della conversazione:", data.message || data);
            }
        },
        async getMyConversation() {
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
            this.messages = []; // <-- svuota subito!
            await this.getConversation(conv.id);
            await this.markMessagesRead();

            // Ogni volta che cambio chat, resetto il polling dei messaggi
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
            const res = await fetch(`${__API_URL__}/conversations/${conversationId}/messages`, {
                headers: { Authorization: userId }
            });
            if (res.ok) {
                const msgs = await res.json();
                
                // Carica le reazioni per ogni messaggio
                for (const msg of msgs) {
                    const reactionRes = await fetch(
                        `${__API_URL__}/conversations/${conversationId}/messages/${msg.id}/reactions`, 
                        { headers: { Authorization: userId } }
                    );
                    msg.reactions = reactionRes.ok ? await reactionRes.json() : [];
                }
                
                this.messages = msgs;
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
                await this.getConversation(this.openConversation.id); // <-- aggiorna subito!
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
        async forwardMessage(targetConversationId) {
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
            // Rimuovi il messaggio inoltrato dalla chat corrente
            this.messages = this.messages.filter(m => m.id !== this.forwardMsg.id);
            this.successMsg = "Messaggio inoltrato!";
            setTimeout(() => { this.successMsg = null; }, 1000); // Nasconde il messaggio dopo 2 secondi
            this.closeForwardModal();
          } else {
            alert("Errore durante l'inoltro del messaggio.");
          }
        },
        async searchUsersGroup() {
          if (this.searchUser.length < 1) {
            this.searchResultsGroup = [];
            this.dropdownOpenGroup = false;
            return;
          }
          const userId = localStorage.getItem("userId");
          const myUsername = localStorage.getItem("username");
          const res = await fetch(`${__API_URL__}/search/users?q=${encodeURIComponent(this.searchUser)}`, {
            headers: { Authorization: userId }
          });
          if (res.ok) {
            const results = await res.json();
            this.searchResultsGroup = results.filter(
              u => u.username !== myUsername && !this.groupMembers.some(m => m.username === u.username)
            );
            this.dropdownOpenGroup = !!this.searchResultsGroup.length;
          } else {
            this.searchResultsGroup = [];
            this.dropdownOpenGroup = false;
          }
        },
        closeDropdownGroup() {
          setTimeout(() => { this.dropdownOpenGroup = false; }, 150);
        },
        addToGroup(user) {
          if (!this.groupMembers.some(u => u.username === user.username)) {
            this.groupMembers.push(user);
          }
          this.searchUser = "";
          this.searchResultsGroup = [];
          this.dropdownOpenGroup = false;
        },
        removeMemberFromGroup(username) {
          this.groupMembers = this.groupMembers.filter(u => u.username !== username);
        },
        async submitGroup() {
          this.groupError = "";
          if (!this.newGroupName || this.newGroupName.length < 3 || this.newGroupName.length > 16) {
            this.groupError = "Il nome deve essere tra 3 e 16 caratteri";
            return;
          }
          if (this.groupMembers.length === 0) {
            this.groupError = "Aggiungi almeno un membro";
            return;
          }
          try {
            const userId = localStorage.getItem("userId");
            const myUsername = localStorage.getItem("username");
            // Prendi tutti gli username selezionati
            let membersArr = this.groupMembers.map(u => u.username);
            // Aggiungi il creatore se non già presente
            if (!membersArr.includes(myUsername)) {
              membersArr.push(myUsername);
            }
            const body = {
              name: this.newGroupName,
              members: membersArr
            };
            if (this.newGroupPhoto) body.photo = this.newGroupPhoto;
            const res = await fetch(`${__API_URL__}/groups`, {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
                Authorization: userId
              },
              body: JSON.stringify(body)
            });
            if (res.ok) {
              this.closeCreateGroupModal();
              await this.getMyConversation?.(); // aggiorna la lista gruppi/conversazioni
            } else {
              const data = await res.json();
              this.groupError = data.message || "Errore creazione gruppo";
            }
          } catch (e) {
            this.groupError = "Errore di rete";
          }
        },
        async submitEditGroup() {
          this.groupError = "";
          if (!this.newGroupName || this.newGroupName.length < 3 || this.newGroupName.length > 16) {
            this.groupError = "Il nome deve essere tra 3 e 16 caratteri";
            return;
          }
          const userId = localStorage.getItem("userId");
          // Cambia nome
          await fetch(`${__API_URL__}/groups/${this.editGroupId}/name`, {
            method: "PATCH",
            headers: { "Content-Type": "application/json", Authorization: userId },
            body: JSON.stringify({ name: this.newGroupName })
          });
          // Cambia foto
          await fetch(`${__API_URL__}/groups/${this.editGroupId}/photo`, {
            method: "PATCH",
            headers: { "Content-Type": "application/json", Authorization: userId },
            body: JSON.stringify({ photo: this.newGroupPhoto })
          });
          // Aggiungi membri (solo quelli nuovi)
          const currentUsernames = (this.openConversation.members || []).map(u => u.username);
          const newMembers = this.groupMembers.filter(u => !currentUsernames.includes(u.username));
          if (newMembers.length > 0) {
            await fetch(`${__API_URL__}/groups/${this.editGroupId}/members`, {
              method: "PATCH",
              headers: { "Content-Type": "application/json", Authorization: userId },
              body: JSON.stringify({ members: newMembers.map(u => u.username) })
            });
          }
          this.closeCreateGroupModal();
          await this.loadAll();
        },
        closeCreateGroupModal() {
          this.openCreateGroupModal = false;
          this.editGroupMode = false;
          this.editGroupId = null;
          this.newGroupName = "";
          this.newGroupPhoto = "";
          this.groupMembers = [];
          this.groupError = "";
          this.searchUser = "";
          this.searchResultsGroup = [];
          this.dropdownOpenGroup = false;
        },
        async loadAll() {
            await this.getMyConversation();
            await this.loadGroups(); 
            if (this.openConversation) {
                await this.getConversation(this.openConversation.id);
                await this.markMessagesRead();
            }
            // Qui puoi aggiungere altre fetch se servono
        },
        async loadGroups() {
            const userId = localStorage.getItem("userId");
            const res = await fetch(`${__API_URL__}/groups`, {
                headers: { Authorization: userId }
            });
            if (res.ok) {
                this.groups = await res.json();
            } else {
                this.groups = [];
            }
        },
        async leaveGroup(groupId) {
            if (!confirm("Sei sicuro di voler lasciare questo gruppo?")) return;
            const userId = localStorage.getItem("userId");
            const res = await fetch(`${__API_URL__}/groups/${groupId}/members`, {
                method: "DELETE",
                headers: { Authorization: userId }
            });
            if (res.ok) {
                this.successMsg = "Hai lasciato il gruppo.";
                await this.loadAll();
                this.openConversation = null;
            } else {
                alert("Errore durante l'uscita dal gruppo.");
            }
        },
        openEditGroupModal() {
            console.log("=== DEBUG EDIT GROUP ===");
            console.log("openConversation:", this.openConversation);
            console.log("openConversation.name:", this.openConversation?.name);
            console.log("openConversation.username:", this.openConversation?.username);
            console.log("openCreateGroupModal prima:", this.openCreateGroupModal);
            
            this.editGroupMode = true;
            this.openCreateGroupModal = true;
            this.editGroupId = this.openConversation.id;
            this.newGroupName = this.openConversation.name || this.openConversation.username;
            this.newGroupPhoto = this.openConversation.photo || this.openConversation.profilePicture || "";
            this.groupMembers = this.openConversation.members || [];
            this.groupError = "";
            
            console.log("openCreateGroupModal dopo:", this.openCreateGroupModal);
            console.log("editGroupMode:", this.editGroupMode);
        },
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