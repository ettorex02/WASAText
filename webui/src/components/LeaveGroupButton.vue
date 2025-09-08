<template>
  <button
    class="btn btn-outline-danger btn-sm ms-3"
    @click="leaveGroup"
    title="Lascia gruppo"
  >
    Lascia Gruppo
  </button>
</template>

<script>
export default {
  name: 'LeaveGroupButton',
  props: {
    groupId: {
      type: [String, Number],
      required: true
    }
  },
  methods: {
    // operationId: leaveGroup
    async leaveGroup() {
      if (!confirm("Sei sicuro di voler lasciare questo gruppo?")) return;
      const userId = localStorage.getItem("userId");
      const res = await fetch(`${__API_URL__}/groups/${this.groupId}/members`, {
        method: "DELETE",
        headers: { Authorization: userId }
      });
      if (res.ok) {
        this.$emit('left');
      } else {
        alert("Errore durante l'uscita dal gruppo.");
      }
    }
  }
}
</script>