<?php
$result = shell_exec('./restarter.sh');
echo '{"Message":"'.$result.'"}'

?>
