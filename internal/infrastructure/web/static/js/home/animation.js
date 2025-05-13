/**
 * Animasi untuk halaman utama Botopia
 */
document.addEventListener('DOMContentLoaded', function() {
    // Animasi fade-in untuk elemen saat scroll
    const fadeElements = document.querySelectorAll('.fade-in');
    
    // Fungsi untuk memeriksa apakah elemen dalam viewport
    function isInViewport(element) {
        const rect = element.getBoundingClientRect();
        return (
            rect.top <= (window.innerHeight || document.documentElement.clientHeight) * 0.85
        );
    }

    // Fungsi untuk mengaktifkan elemen yang terlihat
    function checkFadeElements() {
        fadeElements.forEach(element => {
            if (isInViewport(element)) {
                element.classList.add('active');
            }
        });
    }

    // Periksa saat halaman dimuat dan saat scroll
    checkFadeElements();
    window.addEventListener('scroll', checkFadeElements);

    // Animasi angka pengguna dengan penghitungan
    const counterElements = document.querySelectorAll('.counter-value');
    
    counterElements.forEach(counter => {
        const target = parseInt(counter.getAttribute('data-target'), 10);
        const duration = 1500; // durasi dalam ms
        const step = target / (duration / 16); // 60fps
        
        let current = 0;
        const updateCounter = () => {
            current += step;
            if (current < target) {
                counter.textContent = Math.ceil(current).toLocaleString('id-ID');
                requestAnimationFrame(updateCounter);
            } else {
                counter.textContent = target.toLocaleString('id-ID');
            }
        };
        
        // Mulai animasi hanya ketika elemen dalam viewport
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    updateCounter();
                    observer.unobserve(entry.target);
                }
            });
        }, { threshold: 0.5 });
        
        observer.observe(counter);
    });

    // Animasi rotating untuk testimonial
    const testimonialContainer = document.querySelector('.testimonial-container');
    if (testimonialContainer) {
        const testimonials = testimonialContainer.querySelectorAll('.testimonial-item');
        let currentIndex = 0;

        function showTestimonial(index) {
            testimonials.forEach((testimonial, i) => {
                testimonial.classList.toggle('active', i === index);
            });
        }

        // Tampilkan testimonial pertama
        showTestimonial(currentIndex);

        // Ganti testimonial setiap 5 detik
        setInterval(() => {
            currentIndex = (currentIndex + 1) % testimonials.length;
            showTestimonial(currentIndex);
        }, 5000);
    }
});
