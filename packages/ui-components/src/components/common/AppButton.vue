<template>
    <q-btn :color="color" :label="label" :icon="icon" :loading="loading" :disable="disabled || loading" :dense="dense"
        :size="size" :class="['app-button', `app-button--${variant}`]" v-bind="$attrs" @click="$emit('click', $event)">
        <template v-if="$slots.default" #default>
            <slot />
        </template>
    </q-btn>
</template>

<script setup>
/**
 * AppButton - Shared button component
 *
 * @example
 * <AppButton
 *   label="Submit"
 *   color="primary"
 *   variant="solid"
 *   :loading="isLoading"
 *   @click="handleSubmit"
 * />
 */
defineOptions({
    name: 'AppButton',
    inheritAttrs: false,
})

const props = defineProps({
    label: {
        type: String,
        default: '',
    },
    color: {
        type: String,
        default: 'primary',
        validator: (value) => ['primary', 'secondary', 'positive', 'negative', 'warning', 'info'].includes(value),
    },
    variant: {
        type: String,
        default: 'solid',
        validator: (value) => ['solid', 'outline', 'flat', 'link'].includes(value),
    },
    icon: {
        type: String,
        default: '',
    },
    loading: {
        type: Boolean,
        default: false,
    },
    disabled: {
        type: Boolean,
        default: false,
    },
    dense: {
        type: Boolean,
        default: false,
    },
    size: {
        type: String,
        default: 'md',
        validator: (value) => ['sm', 'md', 'lg', 'xl'].includes(value),
    },
})

defineEmits(['click'])
</script>

<style scoped>
.app-button {
    text-transform: none;
    font-weight: 500;
    letter-spacing: 0;
}

.app-button--outline {
    border: 1px solid currentColor;
    background: transparent !important;
}

.app-button--flat {
    background: transparent !important;
}

.app-button--link {
    background: transparent !important;
    text-decoration: underline;
    padding: 0;
}
</style>