<template>
  <div :class="{ 'loading-animation': true, 'centered': centered }">
    <div class="dot"></div>
    <div class="dot" style="animation-delay: 0.2s"></div>
    <div class="dot" style="animation-delay: 0.4s"></div>
  </div>
</template>
<script>
export default {
    name: 'LoadingAnimation',
    props: {
        centered: {
            type: Boolean,
            default: false
        }
    }
}
</script>

<style lang="scss" scoped>
.loading-animation {
  &:before {
    content: "\200b";
  }

  display: flex;
  align-items: center;
}

.centered {
  justify-content: center;
}

@keyframes fade {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes bounce {
  50% {
    transform: translateY(-3px);
  }
  60% {
    transform: translateY(0px);
  }
}

.dot {
  background-color: black;

  height: 6px;
  width: 6px;
  margin: 2px;
  border-radius: 50%;

  animation-name: fade, bounce;
  animation-duration: 1s, 1s;
  animation-timing-function: ease-in, linear;
  animation-iteration-count: infinite, infinite;
  animation-direction: alternate, normal;
}
</style>
