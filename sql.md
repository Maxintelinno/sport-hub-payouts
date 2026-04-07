#############################
### ตาราง booking_credits ###
#############################

ตารางนี้เก็บ “ยอดคงเหลือเครดิตการจอง” ของ user แต่ละคน
แนะนำให้มี 1 แถวต่อ 1 user

CREATE TABLE public.booking_credits (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    balance numeric(12,2) NOT NULL DEFAULT 0,
    total_earned numeric(12,2) NOT NULL DEFAULT 0,
    total_used numeric(12,2) NOT NULL DEFAULT 0,
    total_expired numeric(12,2) NOT NULL DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),

    CONSTRAINT booking_credits_pkey PRIMARY KEY (id),
    CONSTRAINT booking_credits_user_id_key UNIQUE (user_id),
    CONSTRAINT booking_credits_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);

CREATE INDEX idx_booking_credits_user_id
ON public.booking_credits(user_id);

#########################################
### ตาราง booking_credit_transactions ###
#########################################

อันนี้สำคัญมาก เพราะต้องมี ledger ว่า:
- เครดิตมาจากไหน
- ใช้ไปกับ booking ไหน
- หมดอายุเมื่อไร

CREATE TABLE public.booking_credit_transactions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    booking_credit_id uuid NOT NULL,

    transaction_type varchar(30) NOT NULL,
    -- earn / use / expire / adjust

    source_type varchar(30) NOT NULL,
    -- booking_cancel / booking_payment / admin_adjust / expire_job

    reference_id uuid NULL,
    -- booking_id หรือ id อื่นที่เกี่ยวข้อง

    amount numeric(12,2) NOT NULL,
    balance_before numeric(12,2) NOT NULL DEFAULT 0,
    balance_after numeric(12,2) NOT NULL DEFAULT 0,

    expires_at timestamp NULL,
    note text NULL,

    created_at timestamp NOT NULL DEFAULT now(),

    CONSTRAINT booking_credit_transactions_pkey PRIMARY KEY (id),
    CONSTRAINT booking_credit_transactions_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE,
    CONSTRAINT booking_credit_transactions_booking_credit_id_fkey
        FOREIGN KEY (booking_credit_id) REFERENCES public.booking_credits(id) ON DELETE CASCADE
);

CREATE INDEX idx_booking_credit_transactions_user_id
ON public.booking_credit_transactions(user_id);

CREATE INDEX idx_booking_credit_transactions_reference_id
ON public.booking_credit_transactions(reference_id);

CREATE INDEX idx_booking_credit_transactions_expires_at
ON public.booking_credit_transactions(expires_at);

###############################
### ตาราง owner_settlements ###
###############################

ถ้ายังไม่มี ผมแนะนำ schema นี้เลย
ถ้ามีแล้ว ใช้ส่วน ALTER ด้านล่างแทน

CREATE TABLE public.owner_settlements (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    owner_id uuid NOT NULL,
    booking_id uuid NOT NULL,

    gross_amount numeric(12,2) NOT NULL DEFAULT 0,
    refund_amount numeric(12,2) NOT NULL DEFAULT 0,
    net_revenue numeric(12,2) NOT NULL DEFAULT 0,
    platform_fee numeric(12,2) NOT NULL DEFAULT 0,
    owner_net_amount numeric(12,2) NOT NULL DEFAULT 0,

    status varchar(30) NOT NULL DEFAULT 'pending',
    -- pending / available / paid / reversed

    available_at timestamp NULL,
    paid_at timestamp NULL,
    reversed_at timestamp NULL,

    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),

    CONSTRAINT owner_settlements_pkey PRIMARY KEY (id),
    CONSTRAINT owner_settlements_owner_id_fkey
        FOREIGN KEY (owner_id) REFERENCES public.users(id) ON DELETE CASCADE,
    CONSTRAINT owner_settlements_booking_id_key UNIQUE (booking_id)
);

CREATE INDEX idx_owner_settlements_owner_id
ON public.owner_settlements(owner_id);

CREATE INDEX idx_owner_settlements_status
ON public.owner_settlements(status);

CREATE INDEX idx_owner_settlements_available_at
ON public.owner_settlements(available_at);

เพิ่ม check constraint ที่แนะนำ
กันข้อมูลผิด

###########################
### ความหมายการใช้งานจริง ###
###########################
ตอน user ยกเลิก booking แล้วได้เงินคืนเป็นเครดิต

ระบบจะทำ 3 อย่าง:

A. update bookings
refund_rate = 0.70
refund_amount = 700
refund_method = booking_credit
cancelled_at = now()

B. update payments
refunded_amount = 700
refund_status = partial

C. update booking_credits
balance += 700
total_earned += 700

และ insert booking_credit_transactions

############################################
### ตอน user เอาเครดิตไปใช้จ่าย booking ใหม่ ###
############################################
ระบบจะ:
 - หัก booking_credits.balance
 - เพิ่ม total_used
 - insert transaction type = use

#################################################
### ตอน booking ถูกยกเลิกแล้วต้อง reverse เงิน owner #
#################################################

ระบบจะ update owner_settlements เช่น:
 - gross_amount = 1000
 - refund_amount = 700
 - net_revenue = 300
 - platform_fee = 30
 - owner_net_amount = 270
 - status = pending หรือ reversed

ถ้าคืนเต็ม:
 - status = reversed
 - owner_net_amount = 0

###############################
### transaction_type ที่แนะนำ ###
###############################
ใน booking_credit_transactions.transaction_type
ใช้ได้แบบนี้:
 - earn
 - use
 - expire
 - adjust

########################## 
### source_type ที่แนะนำ ###
########################## 
ใน booking_credit_transactions.source_type

ใช้ได้แบบนี้:
 - booking_cancel
 - booking_payment
 - admin_adjust
 - expire_job

##################################################
### ตัวอย่าง insert ตอนคืนเครดิตจาก cancel booking ###
##################################################
-- สมมติ user มี booking_credits อยู่แล้ว
UPDATE public.booking_credits
SET
    balance = balance + 700,
    total_earned = total_earned + 700,
    updated_at = now()
WHERE user_id = 'USER_UUID';

INSERT INTO public.booking_credit_transactions (
    user_id,
    booking_credit_id,
    transaction_type,
    source_type,
    reference_id,
    amount,
    balance_before,
    balance_after,
    expires_at,
    note
)
VALUES (
    'USER_UUID',
    'BOOKING_CREDIT_UUID',
    'earn',
    'booking_cancel',
    'BOOKING_UUID',
    700,
    0,
    700,
    now() + interval '90 days',
    'Refund from cancelled booking'
);
