<template>
  <section id="contact" class="py-20 bg-gray-50">
    <div class="container mx-auto px-4">
      <h2 class="text-3xl font-bold text-center text-gray-800 mb-12">Hubungi Kami</h2>
      <div class="max-w-2xl mx-auto bg-white p-8 rounded-lg shadow-lg">
        <form>
          <div class="mb-4">
            <label for="name" class="block text-gray-700 font-bold mb-2">Nama</label>
            <input type="text" id="name" class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-600" placeholder="Nama Anda">
          </div>
          <div class="mb-4">
            <label for="email" class="block text-gray-700 font-bold mb-2">Email</label>
            <input type="email" id="email" class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-600" placeholder="email@anda.com">
          </div>
          <div class="mb-4">
            <label for="message" class="block text-gray-700 font-bold mb-2">Pesan</label>
            <textarea id="message" rows="4" class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-600" placeholder="Pesan Anda"></textarea>
          </div>
          <div class="text-center">
            <button type="submit" class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-8 rounded-full transition-colors duration-300">
              Kirim Pesan
            </button>
          </div>
        </form>
      </div>
    </div>
  </section>
</template>

<script setup>
</script>
