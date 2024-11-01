import base from './base.js';

class auth_user extends base {
    constructor() {
        const tableName = 'tb_auth_user';
        const primaryKey = 'id';
        const fields = [
            'id',
            'username',
            'password',
            'level',
            'status',
            'created_at',
            'updated_at',
        ];
        super(tableName, primaryKey, fields);
    }
}


export default auth_user;