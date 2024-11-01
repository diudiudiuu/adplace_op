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
    return api('list').then((res: any) => {
        // 循环遍历数据 for of
        for (const item of res) {
            const children = [];
            for (const project of item.project_list) {
                children.push({
                    id: project.project_id,
                    pid: item.server_id,
                    index: `/project/${project.project_id}`,
                    title: project.project_name,
                });
            }

            children.push({
                id: 'aaaaaaa',
                pid: item.server_id,
                index: '/project_from',
                title: '+项目',
            });

            menus.push({
                id: item.server_id,
                title: item.server_name,
                index: `/project/${item.server_id}`,
                icon: 'Platform',
                children: children
            });
        }

        menus.push({
            id: 'aaaa',
            title: '追加服务器',
            index: '/server_form',
            icon: 'CirclePlusFilled',
        });

        return menus;
    });
}

const menuData = await getMenus()

export {
    menuData
}

