<template>
    <div class="d-flex" style="min-height: 100vh;">
        <!-- Sidebar -->
        <div class="bg-light p-3" style="width: 420px; border-right: 0; box-sizing: border-box;">
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
            const res = await fetch(`http://localhost:3000/search/users?q=${encodeURIComponent(this.search)}`, {
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

            const res = await fetch("http://localhost:3000/conversations", {
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
            } else {
                // Rimuovo l'alert fastidioso, mostro solo in console
                console.log("Errore nella creazione della conversazione:", data.message || data);
            }
        }
    },
    mounted() {
        this.refresh();
        if (this.$route.query.msg) {
            this.successMsg = this.$route.query.msg;
            this.$router.replace({ path: this.$route.path, query: {} });
        }
    }
}
</script>

<style>
</style>
