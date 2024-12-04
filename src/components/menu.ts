import type { Menus } from '@/types/menu';
import api from '@/api';

const menus: Menus[] = [
    // {
    //     id: '1',
    //     title: '服务器1',
    //     index: '/project/1',
    //     icon: 'Platform',
    //     children: [
    //         {
    //             id: '1',
    //             pid: '1',
    //             index: '/project/1',
    //             title: '系统1',
    //         },
    //         {
    //             id: '2',
    //             pid: '1',
    //             index: '/project/2',
    //             title: '系统2',
    //         },
    //     ]
    // }, 
];

const getMenus = () => {
    menus.push({
        id: 'welcome',
        title: '主页',
        index: '/welcome',
        icon: 'Grid',
    });

    // biome-ignore lint/suspicious/noExplicitAny: <explanation>
    return api('list').then((res: any) => {
        // 循环遍历数据 for of
        for (const item of res) {
            const children = [];
            for (const project of item.project_list) {
                children.push({
                    id: project.project_id,
                    pid: item.server_id,
                    index: `/project/${item.server_id}/${project.project_id}`,
                    title: `${project.project_name}`,
                });
            }

            children.push({
                id: `/project_form/${item.server_id}`,
                pid: item.server_id,
                index: `/project_form/${item.server_id}`,
                title: '👉🏼👉添加客户',
            });

            menus.push({
                id: item.server_id,
                title: item.server_name,
                index: `/project/${item.server_id}`,
                icon: 'Platform',
                children: children
            });
        }
        return menus;
    });
}

export {
    getMenus,
    menus
}

