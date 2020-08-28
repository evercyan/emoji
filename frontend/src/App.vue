<template>
  <div id="app" @contextmenu.prevent>
    <Emoji />
  </div>
</template>

<script>
import Emoji from "./components/Emoji.vue";
import "@/assets/css/main.less";
export default {
  name: "app",
  components: {
    Emoji,
  },
};
</script>
