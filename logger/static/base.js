document.addEventListener('alpine:init', () => {


const logList = () => {

  return {
    init() {
      console.log('init')
      this.loadMore()
    },
    items: [], // 存储加载的数据
    page: 1, // 当前页码
    maxPage: 10, // 最大页码
    loading: false, // 是否正在加载
    data: {}, // 当前数据
    // 加载更多数据
    loadMore() {
      if (this.loading || this.page >= this.maxPage) {
        return;
      }

      this.loading = true;

      fetch(`/log/data?page=${this.page}`,
        {
          method: 'GET',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          },
        })
        .then(response => response.json())
        .then(res => {
          this.items = this.items.concat(res.data.map(item => {
            item.time = dayjs.unix(item.time).format('YYYY-MM-DD HH:mm:ss')
            return item
          }));
          this.page += 1;
          this.loading = false;
          console.log(this.items)
        });
    },
    side: {
      ['@scroll'](e) {
        if (e.target.scrollTop + e.target.clientHeight === e.target.scrollHeight) {
          this.loadMore()
        }
      },
    },
  };
};

// 将 Alpine 实例绑定到 DOM 元素上
Alpine.data('loglist', logList);

})