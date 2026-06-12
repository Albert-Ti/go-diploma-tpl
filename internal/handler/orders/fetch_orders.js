const orders = ['9278923470', '2377225624', '12345678903', '346436439', '121160']

let cookie = ''

fetch('http://localhost:8080/api/user/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({login: 'Albert', password: '12345'}),
})
  .then(res => {
    cookie = res.headers.get('set-cookie') || ''
    return res.json()
  })
  .then(() => {
    for (let i = 0; i < orders.length; i++) {
      setTimeout(() => {
        fetch('http://localhost:8080/api/user/orders', {
          method: 'POST',
          headers: {
            'Content-Type': 'text/plain',
            Cookie: cookie,
          },
          body: orders[i],
        })
          .then(async res => {
            const text = await res.text()
            console.log(`Заказ ${orders[i]}:`, res.status, text)
            return {status: res.status, data: text}
          })
          .then(data => console.log(data))
          .catch(err => console.log(err))
      }, 1000 * i)
    }
  })
  .catch(err => console.log(err))
