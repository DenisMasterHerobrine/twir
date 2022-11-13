import { Column, Entity, JoinColumn, ManyToOne, PrimaryColumn, type Relation } from 'typeorm';

import { type User } from './User.js';

@Entity('users_files', { schema: 'public' })
export class UserFile {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'name' })
  name: string;

  @Column('integer', { name: 'size' })
  size: number;

  @Column('text', { name: 'type' })
  type: string;

  @ManyToOne('User', 'files', {
    onDelete: 'SET NULL',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user?: Relation<User>;

  @Column()
  userId: string;
}
