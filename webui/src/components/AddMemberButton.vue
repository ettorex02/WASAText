<template>
  <div class="d-inline">
    <button
      class="btn btn-outline-primary btn-sm me-2"
      title="Aggiungi membri"
      @click="openModal = true"
    >
      + Membri
    </button>
    <div v-if="openModal" class="modal-backdrop">
      <div class="modal-dialog">
        <div class="modal-content p-3">
          <h5>Aggiungi membri al gruppo</h5>
          <input
            v-model="searchUser"
            class="form-control mb-2"
            placeholder="Cerca utenti per username..."
            @input="searchUsers"
            autocomplete="off"
          >
          <ul class="list-group mb-2" style="max-height: 200px; overflow-y: auto;">
            <li
              v-for="user in searchResults"
              :key="user.username"
              class="list-group-item list-group-item-action d-flex align-items-center"
              style="cursor:pointer;"
              @mousedown.prevent="addToGroup(user)"
            >
              <img :src="user.profilePicture" alt="profile" width="32" height="32" class="rounded-circle me-2">
              <span>{{ user.username }}</span>
            </li>
          </ul>
          <button class="btn btn-secondary" @click="closeModal">Chiudi</button>
          <div v-if="addError" class="alert alert-danger mt-2">{{ addError }}</div>
          <div v-if="addSuccess" class="alert alert-success mt-2">{{ addSuccess }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AddMemberButton',
  props: {
    groupId: {
      type: [String, Number],
      required: true
    },
    currentMembers: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      openModal: false,
      searchUser: "",
      searchResults: [],
      addError: "",
      addSuccess: ""
    }
  },
  methods: {
    closeModal() {
      this.openModal = false;
      this.searchUser = "";
      this.searchResults = [];
      this.addError = "";
      this.addSuccess = "";
    },
    async searchUsers() {
      this.addError = "";
      this.addSuccess = "";
      if (this.searchUser.length < 1) {
        this.searchResults = [];
        return;
      }
      const userId = localStorage.getItem("userId");
      const myUsername = localStorage.getItem("username");
      try {
        const res = await this.$axios.get(`/search/users?q=${encodeURIComponent(this.searchUser)}`, {
          headers: { Authorization: userId }
        });
        const results = res.data;
        // Escludi chi è già nel gruppo
        const current = this.currentMembers.map(u => u.username);
        this.searchResults = results.filter(
          u => u.username !== myUsername && !current.includes(u.username)
        );
      } catch {
        this.searchResults = [];
      }
    },
    async addToGroup(user) {
      this.addError = "";
      this.addSuccess = "";
      const userId = localStorage.getItem("userId");
      try {
        await this.$axios.post(`/groups/${this.groupId}/members`, {
          username: user.username
        }, {
          headers: { Authorization: userId }
        });
        this.addSuccess = `Utente ${user.username} aggiunto!`;
        this.searchUser = "";
        this.searchResults = [];
        this.$emit('member-added');
      } catch (e) {
        this.addError = e.response?.data?.message || "Errore aggiunta membro";
      }
    }
  }
}
</script>

<style scoped>
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