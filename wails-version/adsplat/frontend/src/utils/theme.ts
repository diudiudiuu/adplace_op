import { GlobalThemeOverrides } from 'naive-ui'

export const lightThemeOverrides: GlobalThemeOverrides = {
    common: {
        primaryColor: '#3B82F6',
        primaryColorHover: '#2563EB',
        primaryColorPressed: '#1D4ED8',
        primaryColorSuppl: '#60A5FA',
        infoColor: '#0EA5E9',
        successColor: '#10B981',
        warningColor: '#F59E0B',
        errorColor: '#EF4444',
        
        // 基础背景与文字
        baseColor: '#F9FAFB',
        bodyColor: '#FFFFFF',
        textColorBase: '#111827',
        textColor1: '#111827',
        textColor2: '#6B7280',
        textColor3: '#9CA3AF',
        
        // 边框与分割线
        borderColor: '#E5E7EB',
        borderColorHover: '#D1D5DB',
        dividerColor: '#E5E7EB',
        
        // 卡片与悬浮
        cardColor: '#FFFFFF',
        modalColor: '#FFFFFF',
        popoverColor: '#FFFFFF',
        invertedColor: '#1F2937',
        
        // 输入框
        inputColor: '#FFFFFF',
        inputColorDisabled: '#F3F4F6',
        
        // 表格
        tableHeaderColor: '#F9FAFB',
        tableColor: '#FFFFFF',
        tableColorHover: '#F3F4F6',
        tableColorStriped: '#F9FAFB',
        
        // 代码块
        codeColor: '#F3F4F6',
        
        // 滚动条
        scrollbarColor: '#D1D5DB',
        scrollbarColorHover: '#9CA3AF',
    },
    Button: {
        borderRadius: '8px',
        heightMedium: '36px',
        fontWeight: '500',
        colorHoverPrimary: '#2563EB',
        textColorPrimary: '#fff',
        
        // 默认按钮
        color: '#FFFFFF',
        colorHover: '#F9FAFB',
        colorPressed: '#F3F4F6',
        textColor: '#374151',
        textColorHover: '#111827',
        textColorPressed: '#111827',
        border: '1px solid #E5E7EB',
        borderHover: '1px solid #D1D5DB',
        borderPressed: '1px solid #9CA3AF',
    },
    Card: {
        borderRadius: '12px',
        paddingMedium: '20px',
        titleFontSizeMedium: '16px',
        color: '#FFFFFF',
        colorModal: '#FFFFFF',
        colorPopover: '#FFFFFF',
        colorEmbedded: '#F9FAFB',
        textColor: '#111827',
        titleTextColor: '#111827',
        borderColor: '#E5E7EB',
        actionColor: '#F3F4F6',
    },
    Menu: {
        color: '#FFFFFF',
        itemColorHover: '#F3F4F6',
        itemColorActive: '#3B82F6',
        itemColorActiveHover: '#2563EB',
        itemTextColor: '#6B7280',
        itemTextColorHover: '#111827',
        itemTextColorActive: '#FFFFFF',
        itemTextColorActiveHover: '#FFFFFF',
        itemIconColor: '#9CA3AF',
        itemIconColorHover: '#6B7280',
        itemIconColorActive: '#FFFFFF',
        itemIconColorActiveHover: '#FFFFFF',
        arrowColor: '#9CA3AF',
        arrowColorHover: '#6B7280',
        arrowColorActive: '#FFFFFF',
        arrowColorActiveHover: '#FFFFFF',
        borderRadius: '8px',
    },
    Input: {
        color: '#FFFFFF',
        colorFocus: '#FFFFFF',
        colorDisabled: '#F9FAFB',
        textColor: '#111827',
        textColorDisabled: '#9CA3AF',
        placeholderColor: '#9CA3AF',
        placeholderColorDisabled: '#D1D5DB',
        border: '1px solid #E5E7EB',
        borderHover: '1px solid #D1D5DB',
        borderFocus: '1px solid #3B82F6',
        borderDisabled: '1px solid #E5E7EB',
        borderError: '1px solid #EF4444',
        borderWarning: '1px solid #F59E0B',
        borderRadius: '8px',
    },
    DataTable: {
        borderColor: '#E5E7EB',
        thColor: '#F3F4F6',
        tdColor: '#FFFFFF',
        tdColorHover: '#EFF6FF',
        tdColorStriped: '#F9FAFB',
        thTextColor: '#374151',
        tdTextColor: '#111827',
        thFontWeight: '600',
        borderRadius: '8px',
    },
    Layout: {
        color: '#FFFFFF',
        textColor: '#111827',
        siderColor: '#F8FAFC',
        siderBorderColor: '#E2E8F0',
        siderToggleButtonColor: '#F1F5F9',
        siderToggleButtonIconColor: '#64748B',
        siderToggleBarColor: '#E2E8F0',
        headerColor: '#F1F5F9',
        headerBorderColor: '#E2E8F0',
        footerColor: '#F1F5F9',
        footerBorderColor: '#E2E8F0',
    },
    Tabs: {
        colorSegment: '#F3F4F6',
        tabColor: '#FFFFFF',
        tabColorHover: '#F9FAFB',
        tabColorActive: '#3B82F6',
        tabTextColor: '#6B7280',
        tabTextColorHover: '#374151',
        tabTextColorActive: '#FFFFFF',
        tabBorderColor: '#E5E7EB',
        tabBorderColorActive: '#3B82F6',
        paneTextColor: '#111827',
        borderRadius: '8px',
    },
    Modal: {
        color: '#FFFFFF',
        textColor: '#111827',
        titleTextColor: '#111827',
        boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
        borderRadius: '12px',
    },
    Tooltip: {
        color: '#1F2937',
        textColor: '#F9FAFB',
        borderRadius: '6px',
    }
}

export const darkThemeOverrides: GlobalThemeOverrides = {
    common: {
        primaryColor: '#3B82F6',
        primaryColorHover: '#2563EB',
        primaryColorPressed: '#1E40AF',
        infoColor: '#0EA5E9',
        successColor: '#10B981',
        warningColor: '#FBBF24',
        errorColor: '#F87171',
        baseColor: '#111827',
        bodyColor: '#1F2937',
        textColorBase: '#F9FAFB',
        textColor1: '#F9FAFB',
        textColor2: '#D1D5DB',
        textColor3: '#9CA3AF',
        borderColor: '#374151',
        borderColorHover: '#4B5563',
        cardColor: '#1F2937',
        modalColor: '#1F2937',
        popoverColor: '#1F2937',
        invertedColor: '#F9FAFB',
    },
    Card: {
        borderRadius: '12px',
        titleFontSizeMedium: '16px',
    },
    Button: {
        borderRadius: '8px',
        textColorPrimary: '#fff'
    }
}