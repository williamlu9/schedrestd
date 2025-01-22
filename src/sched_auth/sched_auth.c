#include <security/pam_appl.h>
#include <stdlib.h>
#include <string.h>

static struct pam_response *reply = NULL;

int
null_conv(int num_msg, const struct pam_message **msg,
           struct pam_response **resp, void *appdata_ptr)
{
    *resp = reply;
    return PAM_SUCCESS;
}

static struct pam_conv conv = {null_conv, NULL};

// return = true: authenticated
//          false: not authenticated
int
authenticate(char *user, char *password)
{
    pam_handle_t *pamh = NULL;
    int retval;

    if ((retval = pam_start("system-auth", user, &conv, &pamh)) != PAM_SUCCESS)
        return 0;
    reply = calloc(1, sizeof(struct pam_response));
    reply->resp = strdup(password);
    if ((retval = pam_authenticate(pamh, 0)) != PAM_SUCCESS)
       goto auth_end;
    if ((retval = pam_acct_mgmt(pamh, 0)) != PAM_SUCCESS)
       goto auth_end;
    reply = calloc(1, sizeof(struct pam_response));
    reply->resp = strdup(password);
    if ((retval = pam_open_session(pamh, 0)) != PAM_SUCCESS)
       goto auth_end;
    pam_close_session(pamh, 0);
  auth_end:
    pam_end(pamh, PAM_SUCCESS);
    return (retval == PAM_SUCCESS ? 0 : 1);
}

int main(void)
{
    char *user, *pass;
    if ((user = getenv("SCHED_USER")) == NULL ||
        (pass = getenv("SCHED_PASS")) == NULL)
        return 1;
    if (user[0] == 0 || pass[0] == 0)
        return 1;
    return authenticate(user, pass);
}
