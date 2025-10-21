import Base from './base.js';

class UserDomain extends Base {
    constructor() {
        // 表名
        const tableName = 'tb_user_domain';

        // 主键
        const primaryKey = 'id';

        // 显示的字段
        const fields = [
            'id',
            'user_id',
            'domain_id',
            'created_at',
            'updated_at',
        ];

        // 字段类型定义
        const fieldTypes = {
            id: {
                type: 'int',
                value: '',
            },
            user_id: {
                type: 'int',
                value: '',
            },
            domain_id: {
                type: 'int',
                value: '',
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

export default UserDomain;
