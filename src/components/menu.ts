import { Menus } from '@/types/menu';

const menus: Menus[] = [
    {
        id: '1',
        title: '服务器1',
        index: '1',
        icon: 'Platform',
        children: [
            {
                id: '1',
                pid: '1',
                index: '/system1',
                title: '系统1',
            },
            {
                id: '2',
                pid: '1',
                index: '/system2',
                title: '系统2',
            }
        ]
    }, 

    {
        id: '2',
        title: '服务器2',
        index: '2',
        icon: 'Platform',
        children: [
            {
                id: '3',
                pid: '2',
                index: '/system3',
                title: '系统3',
            },
            {
                id: '4',
                pid: '2',
                index: '/system4',
                title: '系统4',
            },
        ],
    }, 

    {
        id: '3',
        title: '追加服务器',
        index: '3',
        icon: 'CirclePlusFilled',
    }, 
];


const getMenus = () => {
    return menus;
}

const menuData = getMenus();

export {
    menuData
}

