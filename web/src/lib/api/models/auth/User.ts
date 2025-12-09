import type Organization from "../app/Organization";

export default interface User {
    first_name: string;
    last_name: string;
    email: string;
    organizations: Organization[];

    updated_at: Date;
    created_at: Date;
}
