import { onMounted, onUnmounted } from 'vue';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';

gsap.registerPlugin(ScrollTrigger);

export function useGsapAnimations() {
  // Stagger fade-in animation for cards
  const animateCards = (selector, options = {}) => {
    const defaults = {
      delay: 0.1,
      stagger: 0.15,
      y: 30,
      duration: 0.8,
      ease: 'power3.out'
    };

    const config = { ...defaults, ...options };

    gsap.fromTo(
      selector,
      {
        opacity: 0,
        y: config.y,
        scale: 0.95
      },
      {
        opacity: 1,
        y: 0,
        scale: 1,
        duration: config.duration,
        delay: config.delay,
        stagger: config.stagger,
        ease: config.ease
      }
    );
  };

  // Scroll-triggered animations
  const animateOnScroll = (selector, options = {}) => {
    const defaults = {
      y: 50,
      duration: 1,
      ease: 'power3.out',
      start: 'top 85%'
    };

    const config = { ...defaults, ...options };

    gsap.fromTo(
      selector,
      {
        opacity: 0,
        y: config.y
      },
      {
        opacity: 1,
        y: 0,
        duration: config.duration,
        ease: config.ease,
        scrollTrigger: {
          trigger: selector,
          start: config.start,
          toggleActions: 'play none none none'
        }
      }
    );
  };

  // Floating animation
  const floatingAnimation = (selector, options = {}) => {
    const defaults = {
      y: -10,
      duration: 2,
      ease: 'sine.inOut'
    };

    const config = { ...defaults, ...options };

    gsap.to(selector, {
      y: config.y,
      duration: config.duration,
      ease: config.ease,
      repeat: -1,
      yoyo: true
    });
  };

  // Page transition
  const pageTransition = (onComplete) => {
    const tl = gsap.timeline({
      onComplete
    });

    tl.to('.page-transition', {
      scaleY: 1,
      duration: 0.5,
      ease: 'power4.inOut',
      transformOrigin: 'bottom'
    }).to('.page-transition', {
      scaleY: 0,
      duration: 0.5,
      ease: 'power4.inOut',
      transformOrigin: 'top'
    });

    return tl;
  };

  // Button hover effect
  const buttonHoverEffect = (selector) => {
    const buttons = document.querySelectorAll(selector);

    buttons.forEach((button) => {
      button.addEventListener('mouseenter', () => {
        gsap.to(button, {
          scale: 1.05,
          duration: 0.3,
          ease: 'power2.out'
        });
      });

      button.addEventListener('mouseleave', () => {
        gsap.to(button, {
          scale: 1,
          duration: 0.3,
          ease: 'power2.out'
        });
      });
    });
  };

  // Number counter animation
  const animateNumber = (element, endValue, duration = 2) => {
    const obj = { value: 0 };

    gsap.to(obj, {
      value: endValue,
      duration,
      ease: 'power1.out',
      onUpdate: () => {
        element.textContent = Math.floor(obj.value);
      }
    });
  };

  // Parallax effect
  const parallaxEffect = (selector, options = {}) => {
    const defaults = {
      yPercent: -30,
      ease: 'none'
    };

    const config = { ...defaults, ...options };

    gsap.to(selector, {
      yPercent: config.yPercent,
      ease: config.ease,
      scrollTrigger: {
        trigger: selector,
        start: 'top bottom',
        end: 'bottom top',
        scrub: 1
      }
    });
  };

  // Reveal text animation
  const revealText = (selector, options = {}) => {
    const defaults = {
      duration: 1,
      stagger: 0.05
    };

    const config = { ...defaults, ...options };

    const elements = document.querySelectorAll(selector);

    elements.forEach((element) => {
      const text = element.textContent;
      element.innerHTML = text
        .split('')
        .map((char) => `<span class="char">${char === ' ' ? '&nbsp;' : char}</span>`)
        .join('');

      gsap.fromTo(
        element.querySelectorAll('.char'),
        {
          opacity: 0,
          y: 20
        },
        {
          opacity: 1,
          y: 0,
          duration: config.duration,
          stagger: config.stagger,
          ease: 'power3.out'
        }
      );
    });
  };

  // Cleanup ScrollTriggers on unmount
  const cleanup = () => {
    ScrollTrigger.getAll().forEach((trigger) => trigger.kill());
  };

  return {
    animateCards,
    animateOnScroll,
    floatingAnimation,
    pageTransition,
    buttonHoverEffect,
    animateNumber,
    parallaxEffect,
    revealText,
    cleanup
  };
}
