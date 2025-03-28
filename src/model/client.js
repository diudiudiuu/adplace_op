import Base from './base.js';

class Client extends Base {
    constructor() {
        // 表名
        const tableName = 'tb_client';

        // 主键
        const primaryKey = 'id';

        // 显示的字段
        const fields = [
            'id',
            'client_id',
            'status',
            'expire_time',
            'api_port',
            'front_port',
            'created_at',
            'updated_at',
        ];

        // 字段类型定义
        const fieldTypes = {
            id: {
                type: 'int',
                value: '',
            },
            client_id: {
                type: 'string',
                value: '',
            },
            status: {
                type: 'enum',
                value: ['normal', 'stop'],
            },
            expire_time: {
                type: 'datetime',
                value: '',
            },
            api_port: {
                type: 'int',
                value: 0,
            },
            front_port: {
                type: 'int',
                value: 0,
            },
            created_at: {
                type: 'datetime',
                value: '',
            },
            updated_at: {
                type: 'datetime',
                value: '',
            }
        };

        super(tableName, primaryKey, fields, fieldTypes);
    }
}

export default Client;
