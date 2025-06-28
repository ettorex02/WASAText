<template>
    <div>
        <div style="position: absolute; top: 80px; right: 30px; z-index: 10;">
            <button class="btn btn-outline-primary me-2" @click="goToProfile">Profilo</button>
            <button class="btn btn-danger" @click="logout">Logout</button>
        </div>

        <div v-if="successMsg" class="alert alert-success text-center" style="max-width: 500px; margin: 0 auto 20px auto;">
            {{ successMsg }}
        </div>

        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
</template>

<script>
export default {
    data: function() {
        return {
            errormsg: null,
            loading: false,
            some_data: null,
            successMsg: null,
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
