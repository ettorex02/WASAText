<template>
  <div class="mt-1 d-flex align-items-center gap-2">
    <!-- Bottone per aggiungere reazione -->
    <button 
      class="btn btn-sm btn-light p-1"
      @click="modalOpen = true"
      title="Aggiungi reazione"
      style="font-size: 1.2rem;"
    >
      😊
    </button>

    <!-- Reazioni accanto al bottone -->
    <div v-if="hasReactions" class="reactions-container ms-2">
      <span
        v-for="(group, emoji) in groupedReactions"
        :key="emoji"
        class="reaction-badge"
        :class="{ 'user-reacted': userHasReacted(group) }"
        :title="getReactionTooltip(group)"
        @click="userHasReacted(group) ? removeReaction() : null"
        style="cursor: pointer;"
      >
        {{ emoji }} <span class="reaction-count">{{ group.length }}</span>
      </span>
    </div>

    <!-- Modal per scegliere emoji -->
    <div v-if="modalOpen" class="modal-backdrop">
      <div class="modal-dialog">
        <div class="modal-content p-3 text-center">
          <h5>Scegli una reazione</h5>
          <div class="d-flex justify-content-center my-3">
            <button
              v-for="emoji in emojis"
              :key="emoji"
              class="btn btn-lg btn-light mx-1 emoji-btn"
              @click="addReaction(emoji)"
            >
              {{ emoji }}
            </button>
          </div>
          <button class="btn btn-secondary mt-2" @click="modalOpen = false">
            Annulla
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: ['message', 'conversationId'],
  data() {
    return {
      modalOpen: false,
      emojis: ['👍', '❤️', '😂', '😮', '😢', '😡'],
      userId: localStorage.getItem("userId"),
    }
  },
  computed: {
    hasReactions() {
      return this.message.reactions && this.message.reactions.length > 0;
    },
    groupedReactions() {
      if (!this.hasReactions) return {};
      const map = {};
      for (const reaction of this.message.reactions) {
        if (!map[reaction.emoji]) {
          map[reaction.emoji] = [];
        }
        map[reaction.emoji].push(reaction);
      }
      return map;
    }
  },
  methods: {
    async addReaction(emoji) {
      try {
        const response = await fetch(
          `${__API_URL__}/conversations/${this.conversationId}/messages/${this.message.id}/reactions`,
          {
            method: "POST",
            headers: { 
              "Content-Type": "application/json", 
              Authorization: this.userId 
            },
            body: JSON.stringify({ emoji }),
          }
        );
        if (response.ok) {
          this.modalOpen = false;
          this.$emit('refresh');
        } else {
          console.error('Errore nell\'aggiunta della reazione');
        }
      } catch (error) {
        console.error('Errore di rete:', error);
      }
    },
    async removeReaction() {
      try {
        const response = await fetch(
          `${__API_URL__}/conversations/${this.conversationId}/messages/${this.message.id}/reactions`,
          {
            method: "DELETE",
            headers: { Authorization: this.userId }
          }
        );
        if (response.ok) {
          this.$emit('refresh');
        } else {
          console.error('Errore nella rimozione della reazione');
        }
      } catch (error) {
        console.error('Errore di rete:', error);
      }
    },
    getReactionTooltip(reactionGroup) {
      // Mostra SOLO username degli utenti che hanno reagito
      return reactionGroup.map(r => r.user.username).join(', ');
    },
    userHasReacted(reactionGroup) {
      return reactionGroup.some(r => String(r.user.id) === this.userId);
    }
  }
}
</script>

<style scoped>
.mt-1.d-flex.align-items-center.gap-2 {
  gap: 0.5rem !important;
}
.reactions-container {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 0;
}
.reaction-badge {
  font-size: 1.1rem;
  background: #e9ecef;
  border-radius: 12px;
  padding: 4px 8px;
  display: inline-block;
  cursor: pointer;
  margin-right: 4px;
  transition: background 0.2s;
}
.reaction-badge.user-reacted {
  background: #d1e7fd;
  color: #0d6efd;
}
.reaction-count {
  font-size: 0.9rem;
  color: #6c757d;
  margin-left: 2px;
}
.emoji-btn {
  font-size: 2rem;
  border: 2px solid transparent;
  transition: all 0.2s;
}
.emoji-btn:hover {
  border-color: #007bff;
  transform: scale(1.1);
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
  width: 90%;
  box-shadow: 0 4px 20px rgba(0,0,0,0.15);
}
</style>