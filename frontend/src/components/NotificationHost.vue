<script setup lang="ts">
import type { AppNotification } from '@/utils';

interface Props {
    notifications: readonly AppNotification[];
}

defineProps<Props>();
</script>

<template>
    <div class="notifications" aria-live="polite" aria-atomic="false">
        <TransitionGroup name="notification">
            <div v-for="notification in notifications" :key="notification.id" class="notification" :class="{ 'notification--wide': notification.type === 'info' }">
                <img v-if="notification.type === 'warning'" src="../assets/icony/alert.svg" alt="Vykřičník" />
                <img v-else-if="notification.type === 'copy'" src="../assets/icony/copy.svg" alt="Zkopírováno" />
                <img v-else src="../assets/icony/info.svg" alt="Oznámení" />
                <span>{{ notification.text }}</span>
            </div>
        </TransitionGroup>
    </div>
</template>

<style scoped>
.notification-move {
    transition: all 0.2s ease;
}

.notification-enter-active,
.notification-leave-active {
    transition: all 0.1s ease;
}

.notification-enter-from,
.notification-leave-to {
    opacity: 0;
    transform: translateX(50px);
}

.notification-leave-active {
    position: absolute;
}

.notifications {
    position: fixed;
    right: 0;
    bottom: 0;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    justify-content: end;
    gap: 10px;
    padding: 20px;
    min-height: 100px;
    pointer-events: none;
    width: 100vw;
    z-index: 1000;
}

.notification {
    min-height: 60px;
    background-color: var(--tmave-fialova);
    min-width: 100px;
    max-width: min(85%, 330px);
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 5px;
    padding: 10px 20px;
    gap: 15px;
    box-shadow: 0 0 10px 2px rgba(0, 0, 0, 0.75);
}

.notification--wide {
    max-width: min(85%, 450px);
}

.notification img {
    width: 24px;
}

.notification span {
    white-space: pre-line;
}
</style>
