<script>
const _ = require('lodash')

export default {
  methods: {
    isCu: function (u) {
      return u.id == this.cuid
    },
    isCuOrAdmin: function (u) {
      return this.isCu(u) || this.isAdmin
    },
  },
  computed: {
    cu: {
      get: function () {
        return this.$root.cu
      },
      set: function (value) {
        this.$root.cu = value
      }
    },
    cuid: function () {
      return _.get(this.cu, 'id', false)
    },
    isAdmin: function() {
      return _.get(this.cu, 'admin', false)
    },
    cuLoading: {
      get: function () {
        return this.$root.cuLoading
      },
      set: function (value) {
        this.$root.cuLoading = value
      }
    }
  }
}
</script>
