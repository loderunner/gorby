<template>
  <nav class="navbar is-fixed-top" role="navigation">
    <div class="navbar-brand">
      <h1 class="navbar-item title">Gorby</h1>
    </div>
    <div class="navbar-menu">
      <div class="navbar-end">
        <div class="navbar-item">
          <div class="control has-icons-left">
            <input id="searchInput" class="input" type="text" placeholder="Search">
            <span class="icon is-small is-left">
              <i class="fas fa-search"></i>
            </span>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<style lang="scss" scoped>
#searchInput {
  width: 300px;
}
</style>
