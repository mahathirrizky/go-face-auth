<template>
  <section class="relative bg-bg-base text-text-base py-24 md:py-32 overflow-hidden min-h-screen flex items-center relative">
    
    <!-- Background Gradient Blobs -->
    <div class="absolute top-0 left-0 w-full h-full overflow-hidden pointer-events-none z-0">
        <div class="absolute top-[-10%] left-[-10%] w-96 h-96 bg-primary/30 rounded-full blur-3xl animate-blob"></div>
        <div class="absolute top-[20%] right-[-10%] w-96 h-96 bg-accent/20 rounded-full blur-3xl animate-blob animation-delay-2000"></div>
        <div class="absolute bottom-[-10%] left-[20%] w-96 h-96 bg-secondary/20 rounded-full blur-3xl animate-blob animation-delay-4000"></div>
    </div>

    <div class="container mx-auto px-4 relative z-10 grid md:grid-cols-2 gap-12 items-center">
      
      <!-- Left Content -->
      <div class="text-left space-y-6">
        <div class="inline-block px-4 py-1.5 rounded-full bg-primary/10 border border-primary/20 text-primary text-sm font-medium mb-4 animate-fade-in-up">
            ✨ Revolusi Absensi Digital
        </div>
        <h1 class="text-5xl md:text-7xl font-bold leading-tight mb-6 animate-fade-in-up animation-delay-200">
          Absensi Wajah <br>
          <span class="text-transparent bg-clip-text bg-gradient-to-r from-primary to-accent">Lebih Akurat.</span>
        </h1>
        <p class="text-lg md:text-xl text-text-muted max-w-lg animate-fade-in-up animation-delay-400 leading-relaxed">
          Tingkatkan efisiensi operasional dengan teknologi pengenalan wajah biometrik canggih. Data hadir real-time, anti-manipulasi.
        </p>
        
        <div class="flex flex-col sm:flex-row gap-4 animate-fade-in-up animation-delay-600 pt-4">
          <Button 
            @click="scrollToPricing" 
            label="Mulai Gratis 14 Hari" 
            class="!bg-primary hover:!bg-primary-dark !border-none !text-white !px-8 !py-3 !rounded-xl !font-semibold !text-lg shadow-lg shadow-primary/30 transition-all hover:scale-105"
          />
          <Button 
            @click="scrollToFeatures" 
            label="Pelajari Cara Kerja" 
            variant="outlined"
            class="!border-white/20 !text-white hover:!bg-white/10 !px-8 !py-3 !rounded-xl !font-semibold !text-lg transition-all"
          />
        </div>

        <div class="pt-8 flex items-center gap-4 text-sm text-text-muted animate-fade-in-up animation-delay-800">
            <div class="flex -space-x-2">
                <div class="w-8 h-8 rounded-full bg-gray-700 border-2 border-bg-base"></div>
                <div class="w-8 h-8 rounded-full bg-gray-600 border-2 border-bg-base"></div>
                <div class="w-8 h-8 rounded-full bg-gray-500 border-2 border-bg-base"></div>
            </div>
            <p>Dipercaya oleh 500+ Perusahaan</p>
        </div>
      </div>

      <!-- Right Visual (Glass Card Mockup) -->
      <div class="relative hidden md:block animate-float">
        <div class="glass rounded-2xl p-6 relative z-10 transform rotate-[-5deg] hover:rotate-0 transition-transform duration-500">
            <!-- Mock Header -->
            <div class="flex items-center justify-between mb-6 border-b border-white/10 pb-4">
                <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-full bg-primary/20 flex items-center justify-center text-primary">
                        <i class="pi pi-user text-xl"></i>
                    </div>
                    <div>
                        <div class="text-sm font-medium">Budi Santoso</div>
                        <div class="text-xs text-text-muted">Marketing Division</div>
                    </div>
                </div>
                <div class="px-3 py-1 rounded-full bg-green-500/20 text-green-400 text-xs">
                    Terverifikasi
                </div>
            </div>
            <!-- Mock Face Scan Visual -->
            <div class="aspect-square bg-black/40 rounded-xl relative overflow-hidden border border-white/5 flex items-center justify-center group">
                <div class="absolute inset-0 bg-gradient-to-b from-transparent via-primary/10 to-transparent opacity-50 animate-scan"></div>
                <i class="pi pi-camera text-4xl text-white/20 group-hover:text-primary transition-colors"></i>
                <div class="absolute bottom-4 left-0 right-0 text-center text-xs text-primary font-mono">
                    Mendeteksi Wajah...
                </div>
            </div>
            <!-- Mock Stats -->
             <div class="grid grid-cols-2 gap-4 mt-6">
                <div class="bg-white/5 p-3 rounded-lg text-center">
                    <div class="text-xs text-text-muted">Masuk</div>
                    <div class="font-bold text-white">08:00 WIB</div>
                </div>
                <div class="bg-white/5 p-3 rounded-lg text-center">
                    <div class="text-xs text-text-muted">Status</div>
                    <div class="font-bold text-green-400">Tepat Waktu</div>
                </div>
             </div>
        </div>
        
        <!-- Decorative Elements behind card -->
        <div class="absolute -top-10 -right-10 w-32 h-32 bg-accent/30 rounded-full blur-2xl"></div>
        <div class="absolute -bottom-10 -left-10 w-40 h-40 bg-primary/30 rounded-full blur-2xl"></div>
      </div>

    </div>
  </section>
</template>

<script setup>
import Button from 'primevue/button';

const scrollToFeatures = () => {
  document.getElementById('features')?.scrollIntoView({ behavior: 'smooth' });
};

const scrollToPricing = () => {
  document.getElementById('pricing')?.scrollIntoView({ behavior: 'smooth' });
};
</script>

<style scoped>
@keyframes blob {
  0% { transform: translate(0px, 0px) scale(1); }
  33% { transform: translate(30px, -50px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
  100% { transform: translate(0px, 0px) scale(1); }
}
.animate-blob {
  animation: blob 7s infinite;
}
.animation-delay-2000 { animation-delay: 2s; }
.animation-delay-4000 { animation-delay: 4s; }

@keyframes scan {
    0% { transform: translateY(-100%); }
    100% { transform: translateY(100%); }
}
.animate-scan {
    animation: scan 2s linear infinite;
}

@keyframes float {
    0%, 100% { transform: translateY(0px); }
    50% { transform: translateY(-20px); }
}
.animate-float {
    animation: float 6s ease-in-out infinite;
}

/* Utils for fade in */
@keyframes fadeInUp {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in-up {
    animation: fadeInUp 0.8s ease-out forwards;
    opacity: 0;
}
.animation-delay-200 { animation-delay: 0.2s; }
.animation-delay-400 { animation-delay: 0.4s; }
.animation-delay-600 { animation-delay: 0.6s; }
.animation-delay-800 { animation-delay: 0.8s; }
</style>