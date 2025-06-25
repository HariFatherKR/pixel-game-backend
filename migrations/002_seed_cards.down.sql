-- Remove all seeded cards
DELETE FROM cards WHERE id IN (
    'code_slash',
    'firewall_up',
    'debug_punch',
    'virus_upload',
    'system_restore',
    'memory_leak',
    'infinite_loop',
    'bug_found',
    'glitch_out',
    'cache_overflow',
    'kernel_panic',
    'basic_attack',
    'basic_defend'
);