/** @type {import('next').NextConfig} */
const withPWA = require('next-pwa')({
  dest: 'public'
})

module.exports = withPWA({
  async redirects() {
    return [
      {
        source: '/',
        destination: '/passages',
        permanent: false,
      },
    ]
  }
})
