import Base from './base.js';

class Domain extends Base {
    constructor() {
        // 表名
        const tableName = 'tb_domain';

        // 主键
        const primaryKey = 'id';

        // 显示的字段
        const fields = [
            'id',
            'domain',
            'expired_at',
            'ping_flag',
            "memo",
            "created_at",
            "updated_at"
        ];

        // 字段类型定义
        const fieldTypes = {
            id: {
                type: 'int',
                value: '',
            },
            domain: {
                type: 'string',
                value: '',
            },
            expired_at: {
                type: 'datetime',
                value: '',
            },
            ping_flag: {
                type: 'enum',
                value: ['deploy', 'error', 'normal'],
            },
            memo: {
                type: 'string',
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

export default Domain;
