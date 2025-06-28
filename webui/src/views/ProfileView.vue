<template>
  <div class="container mt-5" style="max-width:400px">
    <h2>Il tuo profilo</h2>
    <div v-if="user">
      <img :src="user.profilePicture" alt="Profile" width="100" class="mb-3 rounded-circle" />
      <p><b>Username:</b> {{ user.username }}</p>
      <p><b>Display Name:</b> {{ user.displayName }}</p>
    </div>
    <div v-else>
      <p>Caricamento...</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return { user: null }
  },
  async mounted() {
    const username = localStorage.getItem("username")
    const userId = localStorage.getItem("userId")
    const res = await fetch(`http://localhost:3000/users/${username}`, {
      headers: { Authorization: "Bearer " + userId }
    })
    if (res.ok) {
      this.user = await res.json()
    }
  }
}
</script>