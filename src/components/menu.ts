import { Menus } from '@/types/menu';
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
        for (let item of res) {
            const children = [];
            for (let project of item.project_list) {
                children.push({
                    id: project.project_id,
                    pid: item.server_id,
                    index: '/project/' + project.project_id,
                    title: project.project_name,
                });
            }
            menus.push({
                id: item.server_id,
                title: item.server_name,
                index: '/project/' + item.server_id,
                icon: 'Platform',
                children: children
            });
        }

        return menus;
    });
}

const menuData = await getMenus()

export {
    menuData
}

